### 1. dblocker subscribe/publish topic
```
IP address broker: 10.88.81.1
username: DBL0KER
password: 4;1Yf,)`

sub topic: dbl/[serial_numb]/cmd -> untuk menerima perintah
pub topic: dbl/[serial_numb]/sta -> untuk status: offline, online, sleep
pub topic: dbl/[serial_numb]/rpt -> untuk report sensor reading
```
### 2. Live Status dblocker dr mqtt
topic: pub topic: dbl/[serial_numb]/sta

### 3. Menyalakan SSR dengan perintah
topic: dbl/[serial_numb]/cmd

message: dengan 14 sinyal 0 atau 1. 

0 = OFF, 1 = ON

urutan data:
```
|--index--|--yg di switch --|
|    0    |  JAMM-GPS[1]    |
|    1    |  JAM-RC[1]      |
|    2    |  JAMM-GPS[2]    |
|    3    |  JAM-RC[2]      |
|    4    |  JAMM-GPS[3]    |
|    5    |  JAM-RC[3]      |
|    6    |  Kipas Master   |
|    7    |  JAMM-GPS[4]    |
|    8    |  JAM-RC[4]      |
|    9    |  JAMM-GPS[5]    |
|   10    |  JAM-RC[5]      |
|   11    |  JAMM-GPS[6]    |
|   12    |  JAM-RC[6]      |
|   13    |  Kipas Slave    |
```

contoh isi pesan untuk menyalakan semua 14 SSR:

```
1,1,1,1,1,1,1,1,1,1,1,1,1,1
```


server publish dengan setting:
```
h.MqttClient.Publish(topic, 1, true, payload)

qos = 1
retained = true
```

