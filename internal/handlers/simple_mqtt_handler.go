package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	mqtt_client "dblocker_logs_server/internal/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Connect to MQTT
	broker := "tcp://localhost:1883"
	client, err := mqtt_client.NewMqttClient(broker, "simple_handler")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	fmt.Println("Connected to MQTT")

	// 2. Subscribe to a topic
	// We listen to "test/topic" and print whatever comes in
	topic := "test/topic"
	client.Subscribe(topic, func(c mqtt.Client, msg mqtt.Message) {
		fmt.Printf("[MQTT RECEIVED] Topic: %s | Payload: %s\n", msg.Topic(), string(msg.Payload()))
	})

	// 3. Start HTTP Server
	r := gin.Default()

	// POST /send endpoint
	// Reads raw text from body and publishes it to MQTT
	r.POST("/send", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		msg := string(body)

		// Publish to the same topic we are subscribed to
		if err := client.Publish(topic, msg); err != nil {
			c.String(http.StatusInternalServerError, "Error: "+err.Error())
			return
		}

		c.String(http.StatusOK, "Published: "+msg)
	})

	fmt.Println("Server running on :8080. Send POST to /send with text body.")
	r.Run(":8080")
}
