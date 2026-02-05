// Chip: STM32F411CEU6 or STM32F401CCU6, PA11 cmd, WOL PC14 diganti PA15
// MASTER (STM32F411/F401) - HEARTBEAT EDITION
#include <SPI.h>
#include <Ethernet.h>
#include <PubSubClient.h>
#include <max6675.h>
#include <IWatchdog.h> 

#define LED_PIN PC13
#define CMD_PIN PA11 // RS485 DE/RE

// Ethernet
#define W5500_SCK PB3
#define W5500_MISO PB4
#define W5500_MOSI PB5
#define W5500_CS PB1
#define W5500_RST PC15

// Sensors
#define MAX_SCK PB13
#define MAX_MISO PB14
#define MAX_CS_1 PA8  
#define MAX_CS_2 PB2  

HardwareSerial SlaveSerial(PA10, PA9); 

uint32_t outPins[7] = { PB10, PB12, PA12, PB6, PB7, PB8, PB9 };
uint32_t hallSensorPins[9] = { PA0, PA1, PA2, PA3, PA4, PA5, PA6, PA7, PB0 };

MAX6675 temp1(MAX_SCK, MAX_CS_1, MAX_MISO);
MAX6675 temp2(MAX_SCK, MAX_CS_2, MAX_MISO);

byte mac[] = { 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED };
IPAddress mqtt_broker(148, 230, 101, 142); 

char serial_numb[10];
char topic_sub[64];
char topic_pub[64];
char topic_sta[64];

int allHallSensors[18]; 
float temperatures[2];
bool lastSlaveState[7] = { 0 }; 

// Timers
unsigned long lastMqttRetry = 0;
unsigned long lastPublish = 0;
unsigned long lastSlaveMessage = 0;
unsigned long lastHeartbeat = 0; // NEW: Timer for syncing slave

bool slaveConnected = false;
bool isSystemSleeping = false; 

EthernetClient ethClient;
PubSubClient mqttClient(ethClient);

uint8_t crc8(const char* data) {
  uint8_t crc = 0;
  while (*data) { crc ^= (uint8_t)(*data++); }
  return crc;
}

// --- SYNC / HEARTBEAT ---
void syncSlave() {
  char payload[64];
  if (isSystemSleeping) {
     snprintf(payload, sizeof(payload), "SLEEP");
  } else {
     snprintf(payload, sizeof(payload), "SET:%d,%d,%d,%d,%d,%d,%d",
           lastSlaveState[0], lastSlaveState[1], lastSlaveState[2],
           lastSlaveState[3], lastSlaveState[4], lastSlaveState[5], lastSlaveState[6]);
  }

  uint8_t crc = crc8(payload);

  digitalWrite(CMD_PIN, HIGH); delay(2); 
  SlaveSerial.print('$');
  SlaveSerial.print(payload);
  SlaveSerial.print('|');
  if (crc < 0x10) SlaveSerial.print('0');
  SlaveSerial.print(crc, HEX);
  SlaveSerial.println();
  SlaveSerial.flush();
  digitalWrite(CMD_PIN, LOW); 
}

void handleSlaveData(char* rxBuf) {
  lastSlaveMessage = millis();
  slaveConnected = true;

  if (strstr(rxBuf, "REQ:SYNC")) {
    syncSlave(); // Reply instantly to sync request
    return;
  } 
  
  if (strncmp(rxBuf, "CUR:", 4) == 0) {
    // If we are sleeping but Slave is sending data, tell it to SLEEP.
    if (isSystemSleeping) {
        syncSlave();
        return;
    }

    char* ptr = rxBuf + 4;
    for (int i = 9; i < 18; i++) {
      if (ptr) {
        allHallSensors[i] = atoi(ptr);
        ptr = strchr(ptr, ',');
        if (ptr) ptr++;
      }
    }
  }
}

void publishData() {
  if (isSystemSleeping) {
      mqttClient.publish(topic_sta, "SLEEPING");
      return; 
  }
  
  for (int i = 0; i < 9; i++) allHallSensors[i] = analogRead(hallSensorPins[i]);
  float t1 = temp1.readCelsius(); delay(5); 
  float t2 = temp2.readCelsius();
  int t1i = isnan(t1) ? -9900 : (int)(t1 * 100);
  int t2i = isnan(t2) ? -9900 : (int)(t2 * 100);

  if (millis() - lastSlaveMessage > 10000) { // 10s Timeout logic for MQTT reporting
    slaveConnected = false;
    for (int i = 9; i < 18; i++) allHallSensors[i] = 0;
  }

  static char msg[300];
  msg[0] = '\0';
  for (int i = 0; i < 18; i++) {
    char val[12];
    snprintf(val, sizeof(val), "%d,", allHallSensors[i]);
    strcat(msg, val);
  }
  char tail[64];
  snprintf(tail, sizeof(tail), "%d,%d|%d", t1i, t2i, slaveConnected ? 1 : 0);
  strcat(msg, tail);

  mqttClient.publish(topic_pub, msg);
  digitalWrite(LED_PIN, !digitalRead(LED_PIN));
}

void goToSleep() {
  isSystemSleeping = true;
  for (int i = 0; i < 7; i++) digitalWrite(outPins[i], LOW); 
  syncSlave(); 
  mqttClient.publish(topic_sta, "SLEEPING");
}

void mqttCallback(char* topic, byte* payload, unsigned int length) {
  if (length == 5 && strncmp((char*)payload, "SLEEP", 5) == 0) {
    goToSleep();
    return;
  }
  // WAKE_RST: Reset Slave First, then Master
  if (length == 8 && strncmp((char*)payload, "WAKE_RST", 8) == 0) {
    mqttClient.publish(topic_sta, "GLOBAL RESET...");
    digitalWrite(CMD_PIN, HIGH); delay(2);
    SlaveSerial.println("$RESET|");
    SlaveSerial.flush();
    digitalWrite(CMD_PIN, LOW);
    delay(500); 
    NVIC_SystemReset(); 
    return;
  }
  if (length == 9 && strncmp((char*)payload, "RST_SLAVE", 9) == 0) {
    digitalWrite(CMD_PIN, HIGH); delay(2);
    SlaveSerial.println("$RESET|");
    SlaveSerial.flush();
    digitalWrite(CMD_PIN, LOW);
    return;
  }
  if (!isSystemSleeping && length == 2) {
    uint16_t mask = ((uint16_t)payload[0] << 8) | payload[1];
    for (int i = 0; i < 7; i++) digitalWrite(outPins[i], (mask & (1 << i)) ? HIGH : LOW);
    for (int i = 0; i < 7; i++) lastSlaveState[i] = (mask & (1 << (i + 7)));
    syncSlave();
  }
}

void generateIds() {
  snprintf(serial_numb, sizeof(serial_numb), "250001");
  snprintf(topic_sub, sizeof(topic_sub), "dbl/%s/c", serial_numb);
  snprintf(topic_pub, sizeof(topic_pub), "dbl/%s/s", serial_numb);
  snprintf(topic_sta, sizeof(topic_sta), "dbl/%s/sta", serial_numb);
}

void setup() {
  SlaveSerial.begin(9600); 
  pinMode(CMD_PIN, OUTPUT); digitalWrite(CMD_PIN, LOW); 
  pinMode(W5500_RST, OUTPUT); digitalWrite(W5500_RST, LOW); delay(20); digitalWrite(W5500_RST, HIGH); delay(200);
  pinMode(LED_PIN, OUTPUT);

  for (int i = 0; i < 7; i++) {
    pinMode(outPins[i], OUTPUT);
    digitalWrite(outPins[i], LOW);
  }

  generateIds();
  SPI.setMOSI(W5500_MOSI); SPI.setMISO(W5500_MISO); SPI.setSCLK(W5500_SCK); SPI.begin();
  Ethernet.init(W5500_CS);
  Ethernet.begin(mac); 

  mqttClient.setBufferSize(512);
  mqttClient.setServer(mqtt_broker, 1883);
  mqttClient.setCallback(mqttCallback);
  
  // Initial Sync
  syncSlave();

  IWatchdog.begin(15000000); 
}

void loop() {
  IWatchdog.reload(); 
  Ethernet.maintain();

  // --- 1. MASTER HEARTBEAT (The "Keep Alive") ---
  // Send data to slave every 2 seconds regardless of changes
  unsigned long now = millis();
  if (now - lastHeartbeat > 2000) {
      lastHeartbeat = now;
      syncSlave(); // Sends SET or SLEEP, keeps Slave Failsafe happy
  }

  // --- 2. MQTT Logic ---
  if (!mqttClient.connected()) {
    if (now - lastMqttRetry > 5000) {
      lastMqttRetry = now;
      if (mqttClient.connect(serial_numb, NULL, NULL, topic_sta, 1, true, "OFFLINE")) {
        mqttClient.subscribe(topic_sub);
        if(isSystemSleeping) mqttClient.publish(topic_sta, "SLEEPING", true);
        else mqttClient.publish(topic_sta, "ONLINE", true);
        syncSlave();
      }
    }
  } else {
    mqttClient.loop();
  }

  if (now - lastPublish > 2000) {
    lastPublish = now;
    if (mqttClient.connected()) {
      publishData();
    }
  }

  // --- 3. Read Slave Responses ---
  static char rxBuf[128];
  static int rxIdx = 0;
  while (SlaveSerial.available()) {
    char c = SlaveSerial.read();
    if (c == '$') rxIdx = 0; 
    else if (c == '\n' || c == '\r') {
      if (rxIdx > 0) {
        rxBuf[rxIdx] = 0;
        handleSlaveData(rxBuf);
        rxIdx = 0;
      }
    } 
    else if (rxIdx < 127) {
      rxBuf[rxIdx++] = c;
    }
  }
}