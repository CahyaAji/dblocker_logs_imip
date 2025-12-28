package handlers

import (
	"dblocker_logs_server/internal/models"
	"dblocker_logs_server/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	Repo *repository.DeviceRepository
}

func NewDeviceHandler(repo *repository.DeviceRepository) *DeviceHandler {
	return &DeviceHandler{Repo: repo}
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var input models.Device
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func (h *DeviceHandler) GetDevices(c *gin.Context) {
	devices, err := h.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": devices})
}
