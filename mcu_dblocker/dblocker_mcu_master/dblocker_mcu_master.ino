//chip stm32f411ceu6 atau stm32f401ccu6
#include <SPI.h>
#include <Ethernet.h>
#include <PubSubClient.h>

// Temperature
#define MAX_SCK PB13
#define MAX_MISO PB14
#define MAX_MOSI PB15
#define MAX_CS_1 PA8  // Sensor 1
#define MAX_CS_2 PB2  // Sensor 2

// Ethernet
#define W5500_SCK PB3
#define W5500_MISO PB4
#define W5500_MOSI PB5
#define W5500_CS PB1
#define WOL_INT PC14
#define W5500_RST PC15

// Uart ke Slave
#define WAKE_SLAVE PA11
HardwareSerial SlaveSerial(PA10, PA9); //rx, tx

// digital output pin (SSR)
uint32_t jammerPins[7] = { PB10, PB12, PA12, PB6, PB7, PB8, PB9 };

// current sensor (analog read)
uint32_t currentSensorPins[9] = { PA0, PA1, PA2, PA3, PA4, PA5, PA6, PA7, PB0 };
// int currentValue[9];

byte mac[] = { 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED };
IPAddress mqtt_broker(148, 230, 101, 142);

// MQTT, Change size if needed
char serial_numb[10];
char topic_sub[64];
char topic_pub[64];

EthernetClient ethClient;
PubSubClient mqttClient(ethClient);

void generateIds() {
  // Change Serial Number
  snprintf(serial_numb, sizeof(serial_numb), "250001");

  snprintf(topic_sub, sizeof(topic_sub), "dbl/%s/c", serial_numb);
  snprintf(topic_pub, sizeof(topic_pub), "dbl/%s/s", serial_numb);
}

void setup() {
  SlaveSerial.begin(9600);
  SlaveSerial.setTimeout(100);

  pinMode(W5500_RST, OUTPUT);
  digitalWrite(W5500_RST, LOW);
  delay(1);
  digitalWrite(W5500_RST, HIGH);

  pinMode(WAKE_SLAVE, OUTPUT);
  digitalWrite(WAKE_SLAVE, HIGH);

  // Setup output SSR
  for (int i = 0; i < 7; i++) {
    pinMode(jammerPins[i], OUTPUT);
    digitalWrite(jammerPins[i], LOW);
  }

  generateIds();

  // W5500 ethernet
  SPI.setMOSI(W5500_MOSI);
  SPI.setMISO(W5500_MISO);
  SPI.setSCLK(W5500_SCK);
  SPI.begin();

  Ethernet.init(W5500_CS);
  if (Ethernet.begin(mac) != 0) {
    // Debug print
    SlaveSerial.println("Ethernet Connected");
    delay(100);
  } else {
    SlaveSerial.println("error eth");
  }

  SlaveSerial.println("connecting to broker");
  mqttClient.setBufferSize(512);
  mqttClient.setServer(mqtt_broker, 1883);
  mqttClient.setCallback(mqttCallback);

}

void loop() {
  
  Ethernet.maintain();

  if (!mqttClient.connected()) {
    mqttConnect();
  }
  // keep MQTT processing
  mqttClient.loop();
  
}
