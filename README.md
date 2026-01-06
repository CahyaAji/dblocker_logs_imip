## Database Setup

podman run -d --name dblocker-db -e POSTGRES_USER=scm -e POSTGRES_PASSWORD=Menoreh01! -e POSTGRES_DB=dblocker_logs -p 5432:5432 postgres:18.1

## Running the Application

The application is configured via environment variables. To run it safely with your credentials:

```bash
DB_HOST=localhost \
DB_USER=scm \
DB_PASSWORD=Menoreh01! \
DB_NAME=dblocker_logs \
DB_PORT=5432 \
go run .
```

```bash
 DB_PASSWORD=Menoreh01! DB_NAME=dblocker_logs go run ./cmd/api/main.go
```

### mqtt
buat file dulu
mkdir -p ~/mosquitto/config ~/mosquitto/data ~/mosquitto/log

```
~/Files/dev101/mosquitto

podman run -d \
  --name mosquitto \
  -p 1883:1883 \
  -v ~/Files/dev101/mosquitto/config:/mosquitto/config:Z \
  -v ~/Files/dev101/mosquitto/data:/mosquitto/data:Z \
  -v ~/Files/dev101/mosquitto/log:/mosquitto/log:Z \
  eclipse-mosquitto:2.0.22
```