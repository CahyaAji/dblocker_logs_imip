package routes

import (
	"dblocker_logs_server/internal/handlers"
	"dblocker_logs_server/internal/infrastructure/mqtt"
	"dblocker_logs_server/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, mqttClient mqtt.Client) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	deviceRepo := repository.NewDeviceRepository(db)
	dblockerRepo := repository.NewDBlockerRepository(db)
	actionLogRepo := repository.NewActionLogRepository(db)
	deviceLogRepo := repository.NewDeviceLogRepository(db)
	userRepo := repository.NewUserRepository(db)

	deviceControlHandler := handlers.NewDeviceControlHandler(mqttClient, deviceRepo)

	deviceHandler := handlers.NewDeviceHandler(deviceRepo)
	dblockerHandler := handlers.NewDBlockerHandler(dblockerRepo, mqttClient)
	actionLogHandler := handlers.NewActionLogHandler(actionLogRepo)
	deviceLogHandler := handlers.NewDeviceLogHandler(deviceLogRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	// Devices
	api.POST("/devices", deviceHandler.CreateDevice)
	api.GET("/devices", deviceHandler.GetDevices)
	api.GET("/devices/:id", deviceHandler.GetDeviceByID)
	api.PUT("/devices/:id", deviceHandler.UpdateDevice)
	api.DELETE("/devices/:id", deviceHandler.DeleteDevice)

	// DBlockers
	api.POST("/dblockers", dblockerHandler.CreateDBlocker)
	api.GET("/dblockers", dblockerHandler.GetDBlockers)
	api.GET("/dblockers/:id", dblockerHandler.GetDBlockerByID)
	api.PUT("/dblockers/:id", dblockerHandler.UpdateDBlocker)
	api.PUT("/dblockers/config", dblockerHandler.UpdateDBlockerConfig)
	api.DELETE("/dblockers/:id", dblockerHandler.DeleteDBlocker)

	// Commands
	api.POST("/commands", deviceControlHandler.ExecuteCommand)

	// Action logs
	api.POST("/action-logs", actionLogHandler.CreateActionLog)
	api.GET("/action-logs", actionLogHandler.GetActionLogs)
	api.GET("/action-logs/:id", actionLogHandler.GetActionLogByID)
	api.PUT("/action-logs/:id", actionLogHandler.UpdateActionLog)
	api.DELETE("/action-logs/:id", actionLogHandler.DeleteActionLog)

	// Device logs
	api.POST("/device-logs", deviceLogHandler.CreateDeviceLog)
	api.GET("/device-logs", deviceLogHandler.GetDeviceLogs)
	api.GET("/device-logs/:id", deviceLogHandler.GetDeviceLogByID)
	api.PUT("/device-logs/:id", deviceLogHandler.UpdateDeviceLog)
	api.DELETE("/device-logs/:id", deviceLogHandler.DeleteDeviceLog)

	// Users
	api.POST("/users", userHandler.CreateUser)
	api.GET("/users", userHandler.GetUsers)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	return r
}
