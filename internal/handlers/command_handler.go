package handlers

import (
	"dblocker_logs_server/internal/models"
	"dblocker_logs_server/internal/repository"
	mqtt_client "dblocker_logs_server/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

type CommandHandler struct {
	MqttClient *mqtt_client.MqttClient
	DeviceRepo *repository.DeviceRepository
}

func NewCommandHandler(mqttClient *mqtt_client.MqttClient, deviceRepo *repository.DeviceRepository) *CommandHandler {
	return &CommandHandler{
		MqttClient: mqttClient,
		DeviceRepo: deviceRepo,
	}
}

type DeviceCommand struct {
	DeviceID     uint  `json:"device_id"`
	Command      []int `json:"command"`
	WaitResponse bool  `json:"wait_response"`
}

type CommandResult struct {
	DeviceID uint   `json:"device_id"`
	Status   string `json:"status"` // "success" or "error"
	Response string `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}

func (h *CommandHandler) ExecuteCommand(c *gin.Context) {
	var req []DeviceCommand
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "command list is empty"})
		return
	}

	// Execute commands concurrently
	resultsChan := make(chan CommandResult, len(req))
	for _, cmd := range req {
		go func(dc DeviceCommand) {
			result := CommandResult{
				DeviceID: dc.DeviceID,
			}

			device, err := h.DeviceRepo.FindByID(dc.DeviceID)
			if err != nil {
				result.Status = "error"
				result.Error = fmt.Sprintf("device not found: %v", err)
				resultsChan <- result
				return
			}

			var resp string
			if dc.WaitResponse {
				resp, err = h.sendCommandWithResponse(device, dc.Command)
			} else {
				err = h.sendCommandNoResponse(device, dc.Command)
				if err == nil {
					resp = "command sent (no verification)"
				}
			}

			if err != nil {
				result.Status = "error"
				result.Error = err.Error()
			} else {
				result.Status = "success"
				result.Response = resp
			}
			resultsChan <- result
		}(cmd)
	}

	// Collect results
	var results []CommandResult
	for i := 0; i < len(req); i++ {
		results = append(results, <-resultsChan)
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// sendCommandWithResponse sends a command and waits for a response
func (h *CommandHandler) sendCommandWithResponse(device *models.Device, command []int) (string, error) {
	cmdTopic := fmt.Sprintf("devices/%s/cmd", device.SerialNumb)
	respTopic := fmt.Sprintf("devices/%s/response", device.SerialNumb)

	// Channel to receive the response
	// Use a buffered channel to prevent the MQTT callback from blocking (and leaking)
	// if the select times out below.
	responseChan := make(chan string, 1)

	// Subscribe to Response Topic
	msgHandler := func(client mqtt.Client, msg mqtt.Message) {
		responseChan <- string(msg.Payload())
	}

	if err := h.MqttClient.Subscribe(respTopic, msgHandler); err != nil {
		return "", fmt.Errorf("failed to subscribe: %v", err)
	}
	// Cleanup subscription
	defer h.MqttClient.Unsubscribe(respTopic)

	// Publish Command
	payload := map[string]interface{}{
		"command": command,
		"ts":      time.Now().Format(time.RFC3339),
	}

	// The wrapper handles JSON marshaling if we passed the map, but since we have bytes logic elsewhere or
	// want specific marshaling, we can pass bytes. However, the wrapper accepts interface{}.
	// Let's pass the bytes to be safe with existing logic, or pass the map if the wrapper handles it.
	// The wrapper passes interface{} to paho, which handles []byte.
	payloadBytes, _ := json.Marshal(payload)
	if err := h.MqttClient.Publish(cmdTopic, payloadBytes); err != nil {
		return "", fmt.Errorf("failed to publish: %v", err)
	}

	// Wait for Response or Timeout
	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(5 * time.Second):
		return "", errors.New("timeout waiting for device response")
	}
}

// sendCommandNoResponse sends a command without waiting for a response
func (h *CommandHandler) sendCommandNoResponse(device *models.Device, command []int) error {
	cmdTopic := fmt.Sprintf("devices/%s/cmd", device.SerialNumb)

	// Publish Command
	payload := map[string]interface{}{
		"command": command,
		"ts":      time.Now().Format(time.RFC3339),
	}
	payloadBytes, _ := json.Marshal(payload)

	if err := h.MqttClient.Publish(cmdTopic, payloadBytes); err != nil {
		return fmt.Errorf("failed to publish: %v", err)
	}

	return nil
}
