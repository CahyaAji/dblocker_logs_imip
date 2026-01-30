// bool mqttConnect() {
//   // Connect Params: ClientID, Will Topic, Will QoS, Will Retain, Will Message
//   if (mqttClient.connect(serial_numb)) {
//     mqttClient.subscribe(topic_sub);
//     mqttClient.publish(topic_pub, "Connected");
//     notifLed(2); // Success blink
//     return true;
//   }
//   return false;
// }

// // MQTT (Handles incoming messages)
// void mqttCallback(char* topic, byte* payload, unsigned int length) {

//   if (strcmp(topic, topic_sub) == 0) {
//     bool success = decodeBitmask(payload, length);
//     if (!success) {
//       // Debug print
//       // SlaveSerial.println("Gagal decode bitmask, coba CSV...");
//       notifLed(1);
//     }
//   }
// }


// uint8_t crc8(const char *data) {
//   uint8_t crc = 0;
//   while (*data) {
//     crc ^= (uint8_t)(*data++);
//   }
//   return crc;
// }


// void applyOutputs(bool master[7], bool slave[7]) {

//   for (int i = 0; i < 7; i++) {
//     digitalWrite(jammerPins[i], master[i] ? HIGH : LOW);
//   }

//   char slaveCmd[32];
//   strcpy(slaveCmd, "SET:");

//   for (int i = 0; i < 7; i++) {
//     if (i > 0) strcat(slaveCmd, ",");
//     strcat(slaveCmd, slave[i] ? "1" : "0");
//   }

//   uint8_t crc = crc8(slaveCmd);

//   char finalCmd[40];
//   snprintf(finalCmd, sizeof(finalCmd), "%s|%02X", slaveCmd, crc);

//   SlaveSerial.println(finalCmd);
// }



// bool decodeBitmask(byte* payload, unsigned int length) {

//   if (length != 2) return false;

//   uint16_t mask = ((uint16_t)payload[0] << 8) | payload[1];

//   bool master[7] = {
//     mask & (1 << 0),
//     mask & (1 << 1),
//     mask & (1 << 2),
//     mask & (1 << 3),
//     mask & (1 << 4),
//     mask & (1 << 5),
//     mask & (1 << 6),
//   };

//   bool slave[7] = {
//     mask & (1 << 7),
//     mask & (1 << 8),
//     mask & (1 << 9),
//     mask & (1 << 10),
//     mask & (1 << 11),
//     mask & (1 << 12),
//     mask & (1 << 13),
//   };

//   applyOutputs(master, slave);
//   return true;
// }

// bool decodeCSV(byte* payload, unsigned int length) {

//   char buf[64];
//   if (length >= sizeof(buf)) return false;

//   memcpy(buf, payload, length);
//   buf[length] = '\0';

//   bool master[7];
//   bool slave[7];

//   char* token = strtok(buf, ",");
//   int idx = 0;

//   while (token && idx < 14) {
//     bool v = atoi(token) != 0;

//     if (idx < 7) {
//       master[idx] = v;
//     } else {
//       slave[idx - 7] = v;
//     }

//     token = strtok(NULL, ",");
//     idx++;
//   }

//   if (idx != 14) return false;

//   applyOutputs(master, slave);
//   return true;
// }
