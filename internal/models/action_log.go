package models

import (
	"time"

	"gorm.io/datatypes"
)

type ActionLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Timestamp time.Time      `gorm:"autoCreateTime;type:timestamptz(0)" json:"timestamp"`
	UserID    uint           `json:"user_id"`
	UserName  string         `json:"user_name"`
	Action    string         `json:"action"`
	Success   bool           `json:"success"`
	Detail    datatypes.JSON `gorm:"default:null" json:"detail"`
}
