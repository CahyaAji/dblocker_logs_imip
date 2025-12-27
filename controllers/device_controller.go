package controllers

import (
	"dblocker_logs_server/config"
	"dblocker_logs_server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDevice(c *gin.Context) {
	var input models.Device

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	config.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": input})
}

func GetDevices(c *gin.Context) {
	var devices []models.Device
	config.DB.Find((&devices))

	c.JSON(http.StatusOK, gin.H{"success": true, "data": devices})
}
