### 1. Kontroler yg digunakan
- Master stm32f411ceu6
- Slave stm32f401ccu6

### 2. PIN OUT


### 3. LED Indikator
- LED di PC13
```
|-Pola Blink--|--durasi--|----Arti---|
| tidak blink |     0    | normal
| 2x cepat    | 50ms/step| notif sedang berjalan & normal
| 1x lambat   | >=1000ms | notif ada error
```