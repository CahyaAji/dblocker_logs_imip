package models

type Device struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name" binding:"required"`       // e.g., "Living Room Cam"
	Type        string  `json:"type" binding:"required"`       // sensor, actuator, or camera
	IPAddress   string  `json:"ip_address" binding:"required"` // e.g., "192.168.1.50"
	Latitude    float64 `json:"latitude"`                      // e.g., -7.797
	Longitude   float64 `json:"longitude"`                     // e.g., 110.370
	Description string  `json:"description"`                   // Notes about the device
	SerialNumb  string  `json:"serial_numb"`
}
