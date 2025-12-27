package routes

import (
	"dblocker_logs_server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/devices", controllers.CreateDevice)
	r.GET("/devices", controllers.GetDevices)

	return r
}
