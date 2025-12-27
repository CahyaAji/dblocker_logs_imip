package main

import (
	"dblocker_logs_server/config"
	"dblocker_logs_server/routes"
)

func main() {
	config.ConnectDatabase()

	route := routes.SetupRouter()

	// Start server
	route.Run(":5000")
}
