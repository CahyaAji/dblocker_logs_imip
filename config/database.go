package config

import (
	"dblocker_logs_server/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=scm password=Menoreh01! dbname=dblcoker_logs port=5433 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Auto-Migrate the Device model
	database.AutoMigrate(&models.Device{})

	DB = database
	fmt.Println("Database connected successfully!")
}
