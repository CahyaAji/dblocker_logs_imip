// Chip: STM32F411CEU6 or STM32F401CCU6, // PB2 ganti aja ke PB12
// SLAVE (STM32F411/F401) - FAILSAFE EDITION
// #include <Arduino.h>
#include <IWatchdog.h>

#define LED_PIN PC13
#define CMD_PIN PA0 // RS485 Direction Pin (Connect to DE/RE)

// Communication
HardwareSerial CmdSerial(PA10, PA9); 

// Outputs (PB12 used for safety)
uint32_t outPins[7] = { PB10, PB2, PA8, PB6, PB7, PB8, PB9 };

// Sensors
uint32_t hallSensorPins[9] = { PA1, PA2, PA3, PA4, PA5, PA6, PA7, PB0, PB1 };
int allHallSensors[9];

// State Variables
bool isSleeping = false;
unsigned long lastValidPacket = 0; // Timestamp of last Master contact
const unsigned long TIMEOUT_MS = 10000; // 10 Seconds Failsafe

uint8_t crc8(const char* data) {
  uint8_t crc = 0;
  while (*data) { crc ^= (uint8_t)(*data++); }
  return crc;
}

// --- SEND DATA (Only when asked) ---
void replyToMaster() {
  // Read Sensors
  for(int i=0; i<9; i++) allHallSensors[i] = analogRead(hallSensorPins[i]);

  // Enable RS485 TX
  digitalWrite(CMD_PIN, HIGH); delay(2); 

  // Send Packet
  CmdSerial.print("CUR:");
  for(int i=0; i<9; i++) {
    CmdSerial.print(allHallSensors[i]);
    if(i < 8) CmdSerial.print(",");
  }
  CmdSerial.println();
  CmdSerial.flush(); // Wait for data to fly out

  // Disable RS485 TX
  digitalWrite(CMD_PIN, LOW);
}

// --- EMERGENCY SHUTDOWN ---
void failsafeShutdown() {
  // Only shut down if we think we are active
  // If we are already sleeping, we stay sleeping
  if (!isSleeping) {
     for(int i=0; i<7; i++) digitalWrite(outPins[i], LOW);
     // Note: We do NOT set isSleeping=true, because we want to 
     // wake up instantly when the cable reconnects.
  }
}

void processCommand(char* cmd) {
  // Update Heartbeat Timestamp (We are alive!)
  lastValidPacket = millis();

  // 1. SLEEP
  if (strstr(cmd, "SLEEP")) {
    isSleeping = true;
    for(int i=0; i<7; i++) digitalWrite(outPins[i], LOW);
    return;
  }

  // 2. WAKE
  if (strstr(cmd, "WAKE")) {
    isSleeping = false;
    replyToMaster();
    return;
  }

  // 3. RESET
  if (strstr(cmd, "RESET")) {
    delay(100);
    NVIC_SystemReset();
    return;
  }

  // 4. SET (Control & Heartbeat)
  // If we get a SET command, we apply it and reset failsafe.
  if (strncmp(cmd, "SET:", 4) == 0) {
    int states[7];
    int parsed = sscanf(cmd, "SET:%d,%d,%d,%d,%d,%d,%d", 
           &states[0], &states[1], &states[2], &states[3], &states[4], &states[5], &states[6]);
    
    if (parsed == 7) {
      // Force wake if we were sleeping (Auto-Recovery)
      isSleeping = false;

      for(int i=0; i<7; i++) {
        digitalWrite(outPins[i], states[i] ? HIGH : LOW);
      }
      // CRITICAL: Always reply so Master knows we are alive
      replyToMaster();
    }
  }

  // 5. REQ (Manual Ping)
  if (strstr(cmd, "REQ")) {
      lastValidPacket = millis();
      replyToMaster();
  }
}

void verifyAndExecute(char* buf) {
    char* pipePtr = strchr(buf, '|');
    if (!pipePtr) return; 

    *pipePtr = 0; 
    char* payload = buf;
    char* crcHex = pipePtr + 1; 

    if (crc8(payload) == (uint8_t) strtol(crcHex, NULL, 16)) {
        processCommand(payload);
    }
}

void setup(){
  CmdSerial.begin(9600);
  pinMode(CMD_PIN, OUTPUT); digitalWrite(CMD_PIN, LOW); // Listen Mode
  pinMode(LED_PIN, OUTPUT);

  for (int i = 0; i < 7; i++) {
    pinMode(outPins[i], OUTPUT);
    digitalWrite(outPins[i], LOW);
  }
  
  IWatchdog.begin(10000000); // 10s HW Watchdog
  lastValidPacket = millis(); // Initialize timer

  // Ask for Sync on Boot
  digitalWrite(CMD_PIN, HIGH); delay(2);
  CmdSerial.println("REQ:SYNC");
  CmdSerial.flush();
  digitalWrite(CMD_PIN, LOW);
}

void loop(){
  IWatchdog.reload(); 

  // --- 1. RECEIVE DATA ---
  static char rxBuf[64];
  static int rxIdx = 0;

  while (CmdSerial.available()) {
    char c = CmdSerial.read();
    if (c == '$') { rxIdx = 0; digitalWrite(LED_PIN, LOW); } 
    else if (c == '\r' || c == '\n') {
      if (rxIdx > 0) {
        rxBuf[rxIdx] = 0; 
        verifyAndExecute(rxBuf); 
        rxIdx = 0;
      }
      digitalWrite(LED_PIN, HIGH); 
    } 
    else if (rxIdx < 63) { rxBuf[rxIdx++] = c; }
  }

  // --- 2. FAILSAFE CHECK (The Safety Feature) ---
  // If Master hasn't talked to us in 10 seconds, KILL OUTPUTS.
  if (millis() - lastValidPacket > TIMEOUT_MS) {
      failsafeShutdown();
      // Optional: Blink LED fast to indicate error?
  }
}