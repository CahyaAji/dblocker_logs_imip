package database

import (
	"dblocker_logs_server/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() (*gorm.DB, error) {
	dsn := "host=localhost user=scm password=Menoreh01! dbname=dblcoker_logs port=5433 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Device{}, &models.LogEvent{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to Postgres successfully!")
	return db, nil
}
