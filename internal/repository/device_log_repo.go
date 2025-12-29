package repository

import (
	"dblocker_logs_server/internal/models"

	"gorm.io/gorm"
)

type DeviceLogRepository struct {
	DB *gorm.DB
}

func NewDeviceLogRepository(db *gorm.DB) *DeviceLogRepository {
	return &DeviceLogRepository{DB: db}
}

func (r *DeviceLogRepository) Create(log *models.DeviceLog) error {
	return r.DB.Create(log).Error
}

func (r *DeviceLogRepository) FindAll() ([]models.DeviceLog, error) {
	var logs []models.DeviceLog
	err := r.DB.Find(&logs).Error
	return logs, err
}

func (r *DeviceLogRepository) FindByID(id uint) (*models.DeviceLog, error) {
	var log models.DeviceLog
	err := r.DB.First(&log, id).Error
	return &log, err
}

func (r *DeviceLogRepository) Delete(id uint) error {
	return r.DB.Delete(&models.DeviceLog{}, id).Error
}

func (r *DeviceLogRepository) Update(log *models.DeviceLog) error {
	return r.DB.Save(log).Error
}
