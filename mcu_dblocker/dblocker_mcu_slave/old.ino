// -------- READ HANDLER --------
// void handleReadCurrent() {
//   CmdSerial.print("OK:CUR:");
//   for (int i = 0; i < 9; i++) {
//     int v = analogRead(currentSensorPins[i]);
//     CmdSerial.print(v);
//     if (i < 8) CmdSerial.print(',');
//   }
//   CmdSerial.println();
  
//   // debug print
//   Serial.println("Current readings sent via CmdSerial");
// }


// Current sensors (ADC)
// uint32_t currentSensorPins[9] = {
//   PA1, PA2, PA3, PA4, PA5, PA6, PA7, PB0, PB1
// };