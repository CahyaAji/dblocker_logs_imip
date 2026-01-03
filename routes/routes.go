package routes

import (
	"dblocker_logs_server/internal/handlers"
	"dblocker_logs_server/internal/mqtt_client"
	"dblocker_logs_server/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, mqttClient *mqtt_client.MqttClient) *gin.Engine {
	r := gin.Default()

	deviceRepo := repository.NewDeviceRepository(db)
	actionLogRepo := repository.NewActionLogRepository(db)
	deviceLogRepo := repository.NewDeviceLogRepository(db)
	userRepo := repository.NewUserRepository(db)

	commandHandler := handlers.NewCommandHandler(mqttClient, deviceRepo)

	deviceHandler := handlers.NewDeviceHandler(deviceRepo)
	actionLogHandler := handlers.NewActionLogHandler(actionLogRepo)
	deviceLogHandler := handlers.NewDeviceLogHandler(deviceLogRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	r.POST("/devices", deviceHandler.CreateDevice)
	r.GET("/devices", deviceHandler.GetDevices)
	r.GET("/devices/:id", deviceHandler.GetDeviceByID)
	r.PUT("/devices/:id", deviceHandler.UpdateDevice)
	r.DELETE("/devices/:id", deviceHandler.DeleteDevice)

	r.POST("/commands", commandHandler.ExecuteCommand)

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

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	//get user by email

	return r
}