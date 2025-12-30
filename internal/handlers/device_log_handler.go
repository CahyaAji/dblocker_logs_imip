package handlers

import (
	"dblocker_logs_server/internal/models"
	"dblocker_logs_server/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeviceLogHandler struct {
	Repo *repository.DeviceLogRepository
}

func NewDeviceLogHandler(repo *repository.DeviceLogRepository) *DeviceLogHandler {
	return &DeviceLogHandler{Repo: repo}
}

func (h *DeviceLogHandler) CreateDeviceLog(c *gin.Context) {
	var input models.DeviceLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func (h *DeviceLogHandler) GetDeviceLogs(c *gin.Context) {
	deviceLog, err := h.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deviceLog})
}

func (h *DeviceLogHandler) GetDeviceLogByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	deviceLog, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deviceLog})
}

func (h *DeviceLogHandler) UpdateDeviceLog(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input models.DeviceLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = uint(id)
	if err := h.Repo.Update(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

func (h *DeviceLogHandler) DeleteDeviceLog(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.Repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device log deleted successfully"})
}
