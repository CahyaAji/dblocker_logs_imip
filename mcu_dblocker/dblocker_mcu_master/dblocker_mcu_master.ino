//chip stm32f411ceu6 atau stm32f401ccu6
//! hilangkan delay
#include <SPI.h>
#include <Ethernet.h>
#include <PubSubClient.h>

#define LED_PIN PC13

// Ethernet
#define W5500_SCK PB3
#define W5500_MISO PB4
#define W5500_MOSI PB5
#define W5500_CS PB1
#define WOL_INT PC14
#define W5500_RST PC15

// Temperature
#define MAX_SCK PB13
#define MAX_MISO PB14
#define MAX_MOSI PB15
#define MAX_CS_1 PA8  // Sensor 1
#define MAX_CS_2 PB2  // Sensor 2

// Uart ke Slave
#define WAKE_SLAVE PA11
HardwareSerial SlaveSerial(PA10, PA9);  //rx, tx

// digital output pin (SSR)
uint32_t outPins[7] = { PB10, PB12, PA12, PB6, PB7, PB8, PB9 };
//! check SPI2 menggunakan PB12, ganti ke PA15 jika error

// current sensor (analog read)
uint32_t currentSensorPins[9] = { PA0, PA1, PA2, PA3, PA4, PA5, PA6, PA7, PB0 };
// int currentValue[9];

byte mac[] = { 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED };
IPAddress mqtt_broker(148, 230, 101, 142);

// MQTT, Change size if needed
char serial_numb[10];
char topic_sub[64];
char topic_pub[64];

bool lastSlaveState[7] = { 0, 0, 0, 0, 0, 0, 0 };  // Memory of what Slave should be doing
unsigned long lastMqttRetry = 0;                   // For non-blocking timer

EthernetClient ethClient;
PubSubClient mqttClient(ethClient);

uint8_t crc8(const char* data) {
  uint8_t crc = 0;
  while (*data) { crc ^= (uint8_t)(*data++); }
  return crc;
}

void syncSlave() {
  char payload[32];
  // Format: SET:1,0,1,0,1,0,1
  snprintf(payload, sizeof(payload), "SET:%d,%d,%d,%d,%d,%d,%d",
           lastSlaveState[0], lastSlaveState[1], lastSlaveState[2],
           lastSlaveState[3], lastSlaveState[4], lastSlaveState[5], lastSlaveState[6]);

  uint8_t crc = crc8(payload);

  // Wrap in Start ($) and End (\n) markers
  SlaveSerial.print('$');
  SlaveSerial.print(payload);
  SlaveSerial.print('|');
  if (crc < 0x10) SlaveSerial.print('0');  // Hex padding
  SlaveSerial.print(crc, HEX);
  SlaveSerial.println();
}

void mqttCallback(char* topic, byte* payload, unsigned int length) {
  // Check if message is exactly 2 bytes (bitmask)
  if (length == 2) {
    uint16_t mask = ((uint16_t)payload[0] << 8) | payload[1];

    // Update Master Pins (bits 0-6)
    for (int i = 0; i < 7; i++) {
      bool val = mask & (1 << i);
      digitalWrite(outPins[i], val ? HIGH : LOW);
    }

    // Update Slave State Memory (bits 7-13)
    for (int i = 0; i < 7; i++) {
      lastSlaveState[i] = mask & (1 << (i + 7));
    }

    syncSlave();  // Send the new memory to the Slave
  }
}

void generateIds() {
  // Change Serial Number
  snprintf(serial_numb, sizeof(serial_numb), "250001");

  snprintf(topic_sub, sizeof(topic_sub), "dbl/%s/c", serial_numb);
  snprintf(topic_pub, sizeof(topic_pub), "dbl/%s/s", serial_numb);
}

void setup() {
  SlaveSerial.begin(9600);
  //! SlaveSerial.setTimeout(100); apakah ini perlu?

  // Hard Reset W5500
  pinMode(W5500_RST, OUTPUT);
  digitalWrite(W5500_RST, LOW);
  delay(20);
  digitalWrite(W5500_RST, HIGH);
  delay(150);  //delay for ethernet to be ready

  pinMode(LED_PIN, OUTPUT);
  pinMode(WAKE_SLAVE, OUTPUT);
  digitalWrite(WAKE_SLAVE, HIGH);

  // Setup output SSR
  for (int i = 0; i < 7; i++) {
    pinMode(outPins[i], OUTPUT);
    digitalWrite(outPins[i], LOW);
  }

  generateIds();

  // W5500 ethernet
  SPI.setMOSI(W5500_MOSI);
  SPI.setMISO(W5500_MISO);
  SPI.setSCLK(W5500_SCK);
  SPI.begin();

  Ethernet.init(W5500_CS);

  if (Ethernet.begin(mac) == 0) {
    digitalWrite(LED_PIN, LOW);
  } else {
    digitalWrite(LED_PIN, LOW);
    delay(50);
    digitalWrite(LED_PIN, HIGH);
    delay(50);
    digitalWrite(LED_PIN, LOW);
    delay(50);
    digitalWrite(LED_PIN, HIGH);
  }

  // Debug print
  // SlaveSerial.println("connecting to broker");
  // mqttClient.setBufferSize(512);
  mqttClient.setServer(mqtt_broker, 1883);
  mqttClient.setCallback(mqttCallback);
}

void loop() {

  Ethernet.maintain();

  // Non blocking MQTT Reconnect
  if (!mqttClient.connected()) {
    unsigned long now = mills();
    if (now - lastMqttRetry > 5000) {
      lastMqttRetry = now;
      if (mqttClient.connect(serial_numb)) {
        mqttClient.subscribe(topic_sub);
      }
    }
  } else {
    mqttClient.loop();
  }

  if (SlaveSerial.available()) {
    static char rxBuf[32];
    static int rxIdx = 0;
    char c = SlaveSerial.read();

    if (c == '$') rxIdx = 0;
    else if (c == '\n' || c == '\r') {
      rxBuf[rxIdx] = 0;
      if (strstr(rxBuf, "REQ:SYNC")) syncSlave();
      rxIdx = 0;
    } else if (rxIdx < 31) rxBuf[rxIdx++] = c;
  }

  // if (!mqttClient.connected()) {
  //   unsigned long now = millis();
  //   // Attempt to reconnect every 3 seconds without stopping the code
  //   if (now - lastReconnectAttempt > 3000) {
  //     lastReconnectAttempt = now;
  //     if (mqttConnect()) {
  //       lastReconnectAttempt = 0;
  //     } else {
  //       notifLed(1);  // Error blink
  //     }
  //   }
  // } else {
  //   mqttClient.loop();
  // }
}
