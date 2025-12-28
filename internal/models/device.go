package models

type Device struct {
	ID          uint    `gorm:"primaryKey" json:"id"` // Auto-generated ID
	Name        string  `json:"name"`                 // e.g., "Living Room Cam"
	Type        string  `json:"type"`                 // sensor, actuator, or camera
	IPAddress   string  `json:"ip_address"`           // e.g., "192.168.1.50"
	Latitude    float64 `json:"latitude"`             // e.g., -7.797
	Longitude   float64 `json:"longitude"`            // e.g., 110.370
	Description string  `json:"description"`          // Notes about the device
}
