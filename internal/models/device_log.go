package models

import (
	"time"

	"gorm.io/datatypes"
)

type DeviceLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Timestamp  time.Time      `gorm:"autoCreateTime;type:timestamptz(0)" json:"timestamp"`
	DeviceID   uint           `json:"device_id"`
	DeviceName string         `json:"device_name"`
	IsOnline   bool           `json:"is_online"`
	Status     datatypes.JSON `gorm:"default:null" json:"status"`
	Sensors    datatypes.JSON `gorm:"default:null" json:"sensors"`
}
