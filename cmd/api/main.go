package main

import (
	"dblocker_logs_server/internal/database"
	"dblocker_logs_server/internal/infrastructure/mqtt"
	"dblocker_logs_server/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := database.NewPostgresDB()

	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		mqttBroker = "tcp://148.230.101.142:1883"
	}

	mqttClient, err := mqtt.New(mqttBroker, "dblocker-server")
	if err != nil {
		log.Printf("Failed to connect to MQTT broker: %v", err)
	}

	defer mqttClient.Close()

	route := routes.SetupRouter(db, mqttClient)
	// route := routes.SetupRouter(db)

	route.Static("/assets", "./frontend/dist/assets")

	route.GET("/dashboard", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}
	route.Run(":" + port)
}
