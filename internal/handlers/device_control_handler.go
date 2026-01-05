package handlers

import (
	"dblocker_logs_server/internal/infrastructure/mqtt"
	"dblocker_logs_server/internal/models"
	"dblocker_logs_server/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DeviceControlHandler struct {
	MqttClient mqtt.Client
	deviceRepo *repository.DeviceRepository
}

func NewDeviceControlHandler(m mqtt.Client, d *repository.DeviceRepository) *DeviceControlHandler {
	return &DeviceControlHandler{
		MqttClient: m,
		deviceRepo: d,
	}
}

type DeviceCmd struct {
	DeviceID uint  `json:"device_id"`
	Command  []int `json:"command"`
}

type DeviceCmdResponse struct {
	DeviceID uint   `json:"device_id"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
	Response string `json:"response,omitempty"`
}

func (h *DeviceControlHandler) ExecuteCommand(c *gin.Context) {
	var req []DeviceCmd
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "command list is empty"})
		return
	}

	resultsChan := make(chan DeviceCmdResponse, len(req))

	for _, cmd := range req {
		go func(dc DeviceCmd) {
			result := DeviceCmdResponse{DeviceID: dc.DeviceID}

			device, err := h.deviceRepo.FindByID(dc.DeviceID)
			if err != nil {
				result.Status, result.Error = "error", fmt.Sprintf("device not found: %v", err)
				resultsChan <- result
				return
			}

			var resp string
			err = h.sendCmd(device, dc.Command)
			resp = "command published"

			if err != nil {
				result.Status, result.Error = "error", err.Error()
			} else {
				result.Status, result.Response = "success", resp

			}
			resultsChan <- result

		}(cmd)
	}

	var results []DeviceCmdResponse
	for i := 0; i < len(req); i++ {
		results = append(results, <-resultsChan)
	}
	c.JSON(http.StatusOK, results)
}

func (h *DeviceControlHandler) sendCmd(device *models.Device, command []int) error {
	cmdTopic := fmt.Sprintf("dblocker/cmd/%s", device.SerialNumb)
	payload := map[string]any{
		"command": command,
		"ts":      time.Now().Unix(),
	}
	payloadBytes, _ := json.Marshal(payload)

	return h.MqttClient.Publish(cmdTopic, 1, false, payloadBytes)
}

// func (h *CommandHandler) sendCommandWithResponse(device *models.Device, command []int) (string, error) {
// 	cmdTopic := fmt.Sprintf("devices/%s/cmd", device.SerialNumb)
// 	respTopic := fmt.Sprintf("devices/%s/response", device.SerialNumb)

// 	// Buffer of 1 prevents the MQTT goroutine from blocking if the HTTP request times out
// 	responseChan := make(chan string, 1)

// 	// Subscribe to the specific response topic
// 	err := h.mqttClient.Subscribe(respTopic, 1, func(client paho.Client, msg paho.Message) {
// 		responseChan <- string(msg.Payload())
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("subscription failed: %w", err)
// 	}
// 	defer h.mqttClient.Unsubscribe(respTopic)

// 	// Prepare and Publish Payload
// 	payload := map[string]interface{}{
// 		"command": command,
// 		"ts":      time.Now().Unix(), // Using Unix timestamp is often lighter for ESP32
// 	}
// 	payloadBytes, _ := json.Marshal(payload)

// 	if err := h.mqttClient.Publish(cmdTopic, 1, false, payloadBytes); err != nil {
// 		return "", fmt.Errorf("publish failed: %w", err)
// 	}

// 	// Wait for response or timeout
// 	select {
// 	case response := <-responseChan:
// 		return response, nil
// 	case <-time.After(7 * time.Second): // Slightly longer than 5s to account for network lag
// 		return "", errors.New("timeout: device did not respond within 7 seconds")
// 	}
// }
