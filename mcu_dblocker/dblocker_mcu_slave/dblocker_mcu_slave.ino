// Chip: STM32F411CEU6 or STM32F401CCU6
#define LED_PIN PC13

// Master-Slave communication on Pins PA10 (RX) and PA9 (TX)
HardwareSerial CmdSerial(PA10, PA9); 

// Digital outputs (SSR)
uint32_t outPins[7] = { PB10, PB2, PA8, PB6, PB7, PB8, PB9 };

#define CMD_BUF_SIZE 128
char cmdBuf[CMD_BUF_SIZE];
uint8_t cmdIndex = 0;

// -------- CRC-8 --------
uint8_t crc8(const char *data) {
  uint8_t crc = 0;
  while (*data) {
    crc ^= (uint8_t)(*data++);
  }
  return crc;
}

// -------- SET HANDLER --------
void handleSET(char *payload) {
  // debug print
  Serial.print("Parsing SET payload: ");
  Serial.println(payload);

  char *ptr = payload;
  for (int i = 0; i < 7; i++) {
    if (ptr == NULL) {
      // debug print
      Serial.println("Error: Payload ended early!");
      break;
    }

    int val = atoi(ptr);
    digitalWrite(outPins[i], val ? HIGH : LOW);
    
    // debug print
    Serial.print("SSR Index "); Serial.print(i);
    Serial.print(" set to "); Serial.println(val ? "ON" : "OFF");

    ptr = strchr(ptr, ','); 
    if (ptr) ptr++; 
  }
}

// -------- COMMAND PROCESS --------
void processCommand(char *cmd) {
  // debug print
  Serial.print("Raw Command Received: "); Serial.println(cmd);

  char *sep = strchr(cmd, '|');
  if (!sep) {
    // debug print
    Serial.println("Format Error: No CRC separator '|' found");
    return;
  }

  *sep = 0; // Terminate string before CRC
  char *rxCrcStr = sep + 1;

  uint8_t rxCrc = (uint8_t)strtol(rxCrcStr, NULL, 16);
  uint8_t calcCrc = crc8(cmd);

  if (rxCrc != calcCrc) {
    // debug print
    Serial.print("CRC Mismatch! Calculated: "); Serial.print(calcCrc, HEX);
    Serial.print(" Received: "); Serial.println(rxCrc, HEX);
    return;
  }

  if (strncmp(cmd, "SET:", 4) == 0) {
    handleSET(cmd + 4);
    CmdSerial.println("OK:SET");
  } else {
    // debug print
    Serial.print("Unknown Command: "); Serial.println(cmd);
  }
}

// -------- SETUP --------
void setup() {
  // Hardware UART for Master
  CmdSerial.begin(9600);

  // USB CDC Serial for Debugging
  Serial.begin(115200); 
  // Note: On STM32, Serial doesn't need "while(!Serial)" unless you want to wait 
  // for the monitor to open before starting.

  pinMode(LED_PIN, OUTPUT);
  for (int i = 0; i < 7; i++) {
    pinMode(outPins[i], OUTPUT);
    digitalWrite(outPins[i], LOW);
  }

  // debug print
  Serial.println("===============================");
  Serial.println("Slave Awake & Waiting for Master");
  Serial.println("===============================");

  // --- THE HANDSHAKE ---
  delay(1000); // Wait for Master to be ready
  // debug print
  Serial.println("Sending Sync Request to Master...");
  CmdSerial.println("$REQ:SYNC|00"); 
}

// -------- LOOP --------
void loop() {
  while (CmdSerial.available()) {
    char c = CmdSerial.read();

    if (c == '$') { 
      // debug print
      Serial.println("SOF '$' detected, starting new buffer.");
      cmdIndex = 0; 
    } 
    else if (c == '\n' || c == '\r') {
      if (cmdIndex > 0) {
        cmdBuf[cmdIndex] = 0;
        processCommand(cmdBuf);
        cmdIndex = 0;
      }
    } 
    else if (cmdIndex < CMD_BUF_SIZE - 1) {
      cmdBuf[cmdIndex++] = c;
    }
    else {
      // debug print
      Serial.println("Buffer Overflow!");
      cmdIndex = 0;
    }
  }
}