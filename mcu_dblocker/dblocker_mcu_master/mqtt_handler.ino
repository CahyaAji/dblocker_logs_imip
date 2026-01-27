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
      mqttClient.publish(topic_pub, "0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0");
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

  // 1. Safe Buffer Copy
  char msgBuffer[64];
  if (length >= sizeof(msgBuffer)) return;
  memcpy(msgBuffer, payload, length);
  msgBuffer[length] = '\0';
  // Debug print
  // Serial.print("CMD: "); Serial.println(msgBuffer);

  // 2. Parse 14 value
  char slaveCmd[64] = "SET:";
  bool firstSlaveVal = true;
  int idx = 0;
  char* token = strtok(msgBuffer, ",");

  while (token != NULL && idx < 14) {
    int val = atoi(token);

    if (idx < 7) {
      // Local Master Pin
      digitalWrite(jammerPins[idx], val);
    } else {
      // Add to Slave Command
      if (!firstSlaveVal) strcat(slaveCmd, ",");
      strcat(slaveCmd, token);
      firstSlaveVal = false;
    }
    token = strtok(NULL, ",");
    idx++;
  }

  // Send "SET:..." to Slave Immediately
  SlaveSerial.println(slaveCmd);
  // Debug print
  // Serial.print("send to slave-> ");
  // Serial.println(slaveCmd);
}