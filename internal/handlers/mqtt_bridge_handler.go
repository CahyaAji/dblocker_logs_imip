package handlers

import (
	"dblocker_logs_server/internal/infrastructure/mqtt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// MqttBridgeHandler exposes a simple SSE stream that relays MQTT messages
// from the broker to HTTP clients without exposing broker credentials.
type MqttBridgeHandler struct {
	MqttClient mqtt.Client
}

func NewMqttBridgeHandler(mqttClient mqtt.Client) *MqttBridgeHandler {
	return &MqttBridgeHandler{MqttClient: mqttClient}
}

type mqttEvent struct {
	Topic     string `json:"topic"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"ts"`
}

// Stream subscribes to the requested topic and streams messages via SSE.
func (h *MqttBridgeHandler) Stream(c *gin.Context) {
	const topic = "dbl/#"
	const qosVal = byte(0)

	msgCh := make(chan mqttEvent, 16)

	handler := func(msg mqtt.Message) {
		select {
		case msgCh <- mqttEvent{Topic: msg.Topic, Payload: string(msg.Payload), Timestamp: time.Now().Unix()}:
		default:
			// Drop if client is too slow to keep up
		}
	}

	if err := h.MqttClient.Subscribe(topic, qosVal, handler); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe to mqtt topic"})
		return
	}
	defer h.MqttClient.Unsubscribe(topic)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ctx := c.Request.Context()

	c.Stream(func(w io.Writer) bool {
		select {
		case evt := <-msgCh:
			c.SSEvent("mqtt", evt)
			return true
		case <-time.After(30 * time.Second):
			c.SSEvent("ping", "keep-alive")
			return true
		case <-ctx.Done():
			return false
		}
	})
}
