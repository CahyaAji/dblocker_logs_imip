package models

import (
	"time"

	"gorm.io/datatypes"
)

type LogEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
	// Foreign Keys
	DeviceID   uint           `json:"device_id"`         // Must exist
	UserID     *uint          `json:"user_id,omitempty"` // Pointer allows "null" in database
	DeviceName string         `json:"device_name"`
	UserName   string         `json:"user_name,omitempty"`
	Sensors    datatypes.JSON `json:"sensors"`
	Actions    datatypes.JSON `json:"actions"`
}
