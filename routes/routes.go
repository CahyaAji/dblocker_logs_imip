package routes

import (
	"dblocker_logs_server/internal/handlers"
	"dblocker_logs_server/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	deviceRepo := repository.NewDeviceRepository(db)

	deviceHandler := handlers.NewDeviceHandler(deviceRepo)

	r.POST("/devices", deviceHandler.CreateDevice)
	r.GET("/devices", deviceHandler.GetDevices)

	return r
}
