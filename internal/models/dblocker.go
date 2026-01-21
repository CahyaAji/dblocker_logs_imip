package models

type DBlockerConfig struct {
	SignalCtrl bool `json:"signal_ctrl" default:"false"`
	SignalGPS  bool `json:"signal_gps" default:"false"`
}

type DBlocker struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	Name       string           `json:"name" binding:"required"`
	Lat        float64          `json:"latitude" binding:"required"`
	Lng        float64          `json:"longitude" binding:"required"`
	Desc       string           `json:"desc"`
	AngleStart int              `json:"angle_start" default:"0"`
	Config     []DBlockerConfig `gorm:"serializer:json;type:jsonb" json:"config"`
}
