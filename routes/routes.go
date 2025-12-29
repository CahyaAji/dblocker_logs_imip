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
	actionLogRepo := repository.NewActionLogRepository(db)
	deviceLogRepo := repository.NewDeviceLogRepository(db)

	deviceHandler := handlers.NewDeviceHandler(deviceRepo)
	actionLogHandler := handlers.NewActionLogHandler(actionLogRepo)
	deviceLogHandler := handlers.NewDeviceLogHandler(deviceLogRepo)

	r.POST("/devices", deviceHandler.CreateDevice)
	r.GET("/devices", deviceHandler.GetDevices)
	r.GET("/devices/:id", deviceHandler.GetDeviceByID)
	r.PUT("/devices/:id", deviceHandler.UpdateDevice)
	r.DELETE("/devices/:id", deviceHandler.DeleteDevice)

	r.POST("/action-logs", actionLogHandler.CreateActionLog)
	r.GET("/action-logs", actionLogHandler.GetActionLogs)
	r.GET("/action-logs/:id", actionLogHandler.GetActionLogByID)
	r.PUT("/action-logs/:id", actionLogHandler.UpdateActionLog)
	r.DELETE("/action-logs/:id", actionLogHandler.DeleteActionLog)

	r.POST("/device-logs", deviceLogHandler.CreateDeviceLog)
	r.GET("/device-logs", deviceLogHandler.GetDeviceLogs)
	r.GET("/device-logs/:id", deviceLogHandler.GetDeviceLogByID)
	r.PUT("/device-logs/:id", deviceLogHandler.UpdateDeviceLog)
	r.DELETE("/device-logs/:id", deviceLogHandler.DeleteDeviceLog)

	return r
}
