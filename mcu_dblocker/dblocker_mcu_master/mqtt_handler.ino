// //chip stm32f411ceu6 atau stm32f401ccu6
// //! hilangkan delay
// #include <SPI.h>
// #include <Ethernet.h>
// #include <utility/w5100.h>
// #include <PubSubClient.h>
// #include <max6675.h>

// extern "C" void SystemClock_Config(void);

// #define LED_PIN PC13

// // Ethernet
// #define W5500_SCK PB3
// #define W5500_MISO PB4
// #define W5500_MOSI PB5
// #define W5500_CS PB1
// #define W5500_RST PC15
// #define WOL_INT PA15 //PC14

// // Temperature
// #define MAX_SCK PB13
// #define MAX_MISO PB14
// // #define MAX_MOSI PB15
// #define MAX_CS_1 PA8  // Sensor 1
// #define MAX_CS_2 PB2  // Sensor 2

// // Uart ke Slave
// #define WAKE_SLAVE PA11
// HardwareSerial SlaveSerial(PA10, PA9);  //rx, tx

// // digital output pin (SSR)
// uint32_t outPins[7] = { PB10, PB12, PA12, PB6, PB7, PB8, PB9 };

// // current sensor (analog read)
// uint32_t hallSensorPins[9] = { PA0, PA1, PA2, PA3, PA4, PA5, PA6, PA7, PB0 };
// // int currentValue[9];

// MAX6675 temp1(MAX_SCK, MAX_CS_1, MAX_MISO);
// MAX6675 temp2(MAX_SCK, MAX_CS_2, MAX_MISO);

// byte mac[] = { 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED };
// // Ganti dengan IP MQTT Broker
// IPAddress mqtt_broker(148, 230, 101, 142);

// // MQTT, Change size if needed
// char serial_numb[10];
// char topic_sub[64];
// char topic_pub[64];
// char topic_sta[64];

// int allHallSensors[18];
// float temperatures[2];
// bool lastSlaveState[7] = { 0, 0, 0, 0, 0, 0, 0 };

// unsigned long lastMqttRetry = 0;
// unsigned long lastPublish = 0;
// unsigned long lastSlaveMessage = 0;
// bool slaveConnected = false;

// EthernetClient ethClient;
// PubSubClient mqttClient(ethClient);

// uint8_t crc8(const char* data) {
//   uint8_t crc = 0;
//   while (*data) { crc ^= (uint8_t)(*data++); }
//   return crc;
// }

// // --- SYNC TO SLAVE ---
// void syncSlave() {
//   char payload[64];
//   snprintf(payload, sizeof(payload), "SET:%d,%d,%d,%d,%d,%d,%d",
//            lastSlaveState[0], lastSlaveState[1], lastSlaveState[2],
//            lastSlaveState[3], lastSlaveState[4], lastSlaveState[5], lastSlaveState[6]);

//   uint8_t crc = crc8(payload);

//   SlaveSerial.print('$');
//   SlaveSerial.print(payload);
//   SlaveSerial.print('|');
//   if (crc < 0x10) SlaveSerial.print('0');
//   SlaveSerial.print(crc, HEX);
//   SlaveSerial.println();
// }

// // --- HANDLE INCOMING SLAVE DATA ---
// void handleSlaveData(char* rxBuf) {
//   // Update Slave "Heartbeat" timestamp
//   lastSlaveMessage = millis();
//   slaveConnected = true;

//   if (strstr(rxBuf, "REQ:SYNC")) {
//     syncSlave();
//   } else if (strncmp(rxBuf, "CUR:", 4) == 0) {
//     char* ptr = rxBuf + 4;
//     for (int i = 9; i < 18; i++) {
//       if (ptr) {
//         allHallSensors[i] = atoi(ptr);
//         ptr = strchr(ptr, ',');
//         if (ptr) ptr++;
//       }
//     }
//   }
// }


// // --- PUBLISH DATA ---
// void publishData() {
//   // 1. Read Master Hall Sensors
//   for (int i = 0; i < 9; i++) {
//     allHallSensors[i] = analogRead(hallSensorPins[i]);
//   }

//   // 2. Read Temperatures (raw float, but no float formatting)
//   float t1 = temp1.readCelsius();
//   delay(10);
//   float t2 = temp2.readCelsius();

//   // Convert to fixed-point (°C × 100)
//   int t1i = isnan(t1) ? -9900 : (int)(t1 * 100);
//   int t2i = isnan(t2) ? -9900 : (int)(t2 * 100);

//   // 3. Check Slave Health
//   if (millis() - lastSlaveMessage > 5000) {
//     slaveConnected = false;
//     for (int i = 9; i < 18; i++) {
//       allHallSensors[i] = 0;
//     }
//   }

//   // 4. Build MQTT Message
//   static char msg[300];
//   msg[0] = '\0';

//   // 18 Hall Sensors
//   for (int i = 0; i < 18; i++) {
//     char val[12];
//     snprintf(val, sizeof(val), "%d,", allHallSensors[i]);
//     strcat(msg, val);
//   }

//   // Temps (fixed-point) + slave flag
//   // Format:
//   // H0,H1,...,H17,T1i,T2i|S
//   char tail[64];
//   snprintf(tail, sizeof(tail), "%d,%d|%d",
//            t1i, t2i,
//            slaveConnected ? 1 : 0);

//   strcat(msg, tail);

//   // 5. Publish
//   mqttClient.publish(topic_pub, msg);

//   // Activity indicator
//   digitalWrite(LED_PIN, !digitalRead(LED_PIN));
// }


// // --- MQTT CALLBACK ---
// void mqttCallback(char* topic, byte* payload, unsigned int length) {

//   if (length == 5 && strncmp((char*)payload, "SLEEP", 5) == 0) {
//     goToSleep();
//     return;
//   }

//   if (length == 2) {
//     uint16_t mask = ((uint16_t)payload[0] << 8) | payload[1];

//     for (int i = 0; i < 7; i++) {
//       digitalWrite(outPins[i], (mask & (1 << i)) ? HIGH : LOW);
//     }
//     for (int i = 0; i < 7; i++) {
//       lastSlaveState[i] = (mask & (1 << (i + 7)));
//     }
//     syncSlave();
//   }
// }

// void generateIds() {
//   snprintf(serial_numb, sizeof(serial_numb), "250001");
//   snprintf(topic_sub, sizeof(topic_sub), "dbl/%s/c", serial_numb);
//   snprintf(topic_pub, sizeof(topic_pub), "dbl/%s/s", serial_numb);
//   snprintf(topic_sta, sizeof(topic_sta), "dbl/%s/sta", serial_numb); // Status Topic
// }

// void setup() {
//   SlaveSerial.begin(9600);
//   //! SlaveSerial.setTimeout(100); apakah ini perlu?

//   // Hard Reset W5500
//   pinMode(W5500_RST, OUTPUT);
//   digitalWrite(W5500_RST, LOW);
//   delay(20);
//   digitalWrite(W5500_RST, HIGH);
//   delay(200);  //delay for ethernet to be ready

//   pinMode(LED_PIN, OUTPUT);
//   pinMode(WAKE_SLAVE, OUTPUT);
//   digitalWrite(WAKE_SLAVE, HIGH);

//   pinMode(WOL_INT, INPUT_PULLUP);

//   // Setup output SSR
//   for (int i = 0; i < 7; i++) {
//     pinMode(outPins[i], OUTPUT);
//     digitalWrite(outPins[i], LOW);
//   }

//   generateIds();

//   __HAL_RCC_PWR_CLK_ENABLE();

//   // W5500 ethernet
//   SPI.setMOSI(W5500_MOSI);
//   SPI.setMISO(W5500_MISO);
//   SPI.setSCLK(W5500_SCK);
//   SPI.begin();

//   Ethernet.init(W5500_CS);

//   if (Ethernet.begin(mac) == 0) {
//     digitalWrite(LED_PIN, LOW);
//   } else {
//     digitalWrite(LED_PIN, LOW);
//     delay(50);
//     digitalWrite(LED_PIN, HIGH);
//     delay(50);
//     digitalWrite(LED_PIN, LOW);
//     delay(50);
//     digitalWrite(LED_PIN, HIGH);
//   }

//   // Debug print
//   // SlaveSerial.println("connecting to broker");
//   mqttClient.setBufferSize(512);
//   mqttClient.setServer(mqtt_broker, 1883);
//   mqttClient.setCallback(mqttCallback);
// }

// void loop() {
//   Ethernet.maintain();

//   // 1. CONNECTION MANAGER
//   if (!mqttClient.connected()) {
//     unsigned long now = millis();
//     if (now - lastMqttRetry > 5000) {
//       lastMqttRetry = now;
      
//       // Connect with LWT (Last Will and Testament)
//       // If we die, broker publishes "OFFLINE" to topic_sta
//       if (mqttClient.connect(serial_numb, NULL, NULL, topic_sta, 1, true, "OFFLINE")) {
//         mqttClient.subscribe(topic_sub);
//         // We are alive, publish ONLINE
//         mqttClient.publish(topic_sta, "ONLINE", true);
//       }
//     }
//   } else {
//     mqttClient.loop();
//   }

//   // 2. PUBLISH TIMER (Run even if disconnected to keep logic moving)
//   unsigned long now = millis();
//   if (now - lastPublish > 2000) {
//     lastPublish = now;
//     if (mqttClient.connected()) {
//       publishData();
//     }
//   }

//   // 3. SLAVE LISTENER
//   static char rxBuf[128];
//   static int rxIdx = 0;
//   while (SlaveSerial.available()) {
//     char c = SlaveSerial.read();
//     if (c == '$') rxIdx = 0;
//     else if (c == '\n' || c == '\r') {
//       if (rxIdx > 0) {
//         rxBuf[rxIdx] = 0;
//         handleSlaveData(rxBuf);
//         rxIdx = 0;
//       }
//     } else if (rxIdx < 127) {
//       rxBuf[rxIdx++] = c;
//     }
//   }
// }



// void wakeUpHandler() {
//   // Interrupt handler must exist, but does nothing
// }

// void enableW5500WOL() {
//   // Enable WOL bit in MR (Bit 5)
//   uint8_t mr = W5100.readMR();
//   W5100.writeMR(mr | 0x20); 
  
//   // Enable Interrupt Mask for Magic Packet (Bit 5)
//   uint8_t imr = W5100.readIMR();
//   W5100.writeIMR(imr | 0x20); 
// }

// void goToSleep() {
//   mqttClient.publish(topic_sta, "SLEEPING");
//   delay(100); 

//   // 1. Setup WOL on W5500
//   enableW5500WOL();

//   // 2. Attach Interrupt to PC14 (WOL_INT)
//   // W5500 INT pin goes LOW (Falling) on packet detection
//   attachInterrupt(digitalPinToInterrupt(WOL_INT), wakeUpHandler, FALLING);
  
//   // 3. Suspend Tick (Stop internal timers)
//   HAL_SuspendTick();
  
//   // 4. Enter STOP Mode (Deep Sleep)
//   HAL_PWR_EnterSTOPMode(PWR_LOWPOWERREGULATOR_ON, PWR_STOPENTRY_WFI);
  
//   // --- DEVICE SLEEPS HERE ---
  
//   // Restore Clock (Essential!)
//   SystemClock_Config();
//   // 5. Wake Up Routine
//   HAL_ResumeTick(); // Restore timers

//   for(int i=0; i<5; i++){
//     digitalWrite(LED_PIN, !digitalRead(LED_PIN));
//     delay(100);
//   }

//   detachInterrupt(digitalPinToInterrupt(WOL_INT));
  
//   // Disable WOL on W5500 so normal interrupts work again
//   uint8_t mr = W5100.readMR();
//   W5100.writeMR(mr & ~0x20); 

//   mqttClient.publish(topic_sta, "ONLINE");
//   lastMqttRetry = 0; 
// }