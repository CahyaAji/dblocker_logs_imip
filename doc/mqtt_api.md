### 1. dblocker subscribe/publish topic
```
sub topic: dbl/[serial_numb]/c
pub topic: dbl/[serial_numb]/s
```

### 2. Menyalakan SSR dengan perintah
topic: dbl/[serial_numb]/c

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