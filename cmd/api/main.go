package main

import (
	"dblocker_logs_server/internal/database"
	"dblocker_logs_server/routes"
	"log"
	"os"
)

func main() {

	db, err := database.NewPostgresDB()

	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	route := routes.SetupRouter(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}
	route.Run(":" + port)
}