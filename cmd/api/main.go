package main

import (
	"dblocker_logs_server/internal/database"
	"dblocker_logs_server/routes"
	"log"
)

func main() {

	db, err := database.NewPostgresDB()

	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	route := routes.SetupRouter(db)
	route.Run(":5000")
}
