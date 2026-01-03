package main

import (
	"dblocker_logs_server/internal/database"
	"dblocker_logs_server/internal/mqtt_client"
	"dblocker_logs_server/routes"
	"log"
	"os"
)

func main() {

	db, err := database.NewPostgresDB()

	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		mqttBroker = "tcp://localhost:1883"
	}

	mqttClient, err := mqtt_client.NewMqttClient(mqttBroker, "dblocker-server")
	if err != nil {
		log.Printf("Failed to connect to MQTT broker: %v", err)
	}

	route := routes.SetupRouter(db, mqttClient)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}
	route.Run(":" + port)
}
