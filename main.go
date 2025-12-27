package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Device struct {
	ID          uint    `gorm:"primaryKey" json:"id"` // Auto-generated ID
	Name        string  `json:"name"`                 // e.g., "Living Room Cam"
	Type        string  `json:"type"`                 // sensor, actuator, or camera
	IPAddress   string  `json:"ip_address"`           // e.g., "192.168.1.50"
	Latitude    float64 `json:"latitude"`             // e.g., -7.797
	Longitude   float64 `json:"longitude"`            // e.g., 110.370
	Description string  `json:"description"`          // Notes about the device
}

var db *gorm.DB

func main() {
	dsn := "host=localhost user=scm password=Menoreh01! dbname=dblcoker_logs port=5433 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. Check port and password! Error:", err)
	}

	// This magically creates the table if it doesn't exist
	db.AutoMigrate(&Device{})
	log.Println("Database connected and migrated successfully!")

	r := gin.Default()
	r.POST("/devices", createDevice)
	r.GET("devices", getDevices)

	log.Println("server starting on http://localhost:5000")
	r.Run(":5000")

}

// Create a new device entry
func createDevice(c *gin.Context) {
	var input Device

	// Read JSON body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to DB
	result := db.Create(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device registered!", "data": input})
}

// Get all devices
func getDevices(c *gin.Context) {
	var devices []Device
	db.Find(&devices)
	c.JSON(http.StatusOK, devices)
}
