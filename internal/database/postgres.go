package database

import (
	"dblocker_logs_server/internal/models"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "scm"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "Menoreh01!"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "dblocker_logs"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	timezone := os.Getenv("DB_TIMEZONE")
	if timezone == "" {
		timezone = "UTC"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Device{}, &models.DeviceLog{}, &models.ActionLog{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to Postgres successfully!")
	return db, nil
}
