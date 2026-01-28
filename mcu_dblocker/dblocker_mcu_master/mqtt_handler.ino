// MQTT connect to broker
void mqttConnect() {
  while (!mqttClient.connected()) {
    // Debug print
    SlaveSerial.println("Attempting MQTT connection to 148.230.101.142...");

    // Connect Params:
    // 1. ClientID: serial_numb (Must be unique!)
    // 2. User/Pass: NULL (Add if needed)
    // 3. Will Topic: topic_pub ("device/000001/sta")
    // 4. Will QoS: 0 no guarantee, 1 (At least once), or 2 (Exactly once)
    // 5. Will Retain: true
    // 6. Will Message: "OFFLINE"
    // if (mqttClient.connect(serial_numb, NULL, NULL, topic_pub, 0, true, "OFFLINE")) {
    if (mqttClient.connect(serial_numb)) {
      mqttClient.subscribe(topic_sub);

      // Debug print
      SlaveSerial.print("Subscribed: ");
      SlaveSerial.println(topic_sub);
      mqttClient.publish(topic_pub, "ini coba mqttConnect");
    } else {
      // Debug print
      SlaveSerial.print("failed, rc=");
      SlaveSerial.print(mqttClient.state());
      SlaveSerial.println(" retry in 5s");
      delay(5000);
      // coba reset W5500 nya kemudian connect, atau reset system
    }
  }
}

// MQTT (Handles incoming messages)
void mqttCallback(char* topic, byte* payload, unsigned int length) {
  
  if (strcmp(topic, topic_sub) == 0) {
    bool success = decodeBitmask(payload, length);
    if (!success) {
      // Debug print
      SlaveSerial.println("Gagal decode bitmask, coba CSV...");
    }
  }
}


void applyOutputs(bool master[7], bool slave[7]) {

  // Master pins
  for (int i = 0; i < 7; i++) {
    digitalWrite(jammerPins[i], master[i] ? HIGH : LOW);
  }

  // Slave command
  char slaveCmd[32];
  strcpy(slaveCmd, "SET:");

  for (int i = 0; i < 7; i++) {
    if (i > 0) strcat(slaveCmd, ",");
    strcat(slaveCmd, slave[i] ? "1" : "0");
  }

  SlaveSerial.println(slaveCmd);
}


bool decodeBitmask(byte* payload, unsigned int length) {

  if (length != 2) return false;

  uint16_t mask = ((uint16_t)payload[0] << 8) | payload[1];

  bool master[7] = {
    mask & (1 << 0),
    mask & (1 << 1),
    mask & (1 << 2),
    mask & (1 << 3),
    mask & (1 << 4),
    mask & (1 << 5),
    mask & (1 << 6),
  };

  bool slave[7] = {
    mask & (1 << 7),
    mask & (1 << 8),
    mask & (1 << 9),
    mask & (1 << 10),
    mask & (1 << 11),
    mask & (1 << 12),
    mask & (1 << 13),
  };

  applyOutputs(master, slave);
  return true;
}

bool decodeCSV(byte* payload, unsigned int length) {

  char buf[64];
  if (length >= sizeof(buf)) return false;

  memcpy(buf, payload, length);
  buf[length] = '\0';

  bool master[7];
  bool slave[7];

  char* token = strtok(buf, ",");
  int idx = 0;

  while (token && idx < 14) {
    bool v = atoi(token) != 0;

    if (idx < 7) {
      master[idx] = v;
    } else {
      slave[idx - 7] = v;
    }

    token = strtok(NULL, ",");
    idx++;
  }

  if (idx != 14) return false;

  applyOutputs(master, slave);
  return true;
}
