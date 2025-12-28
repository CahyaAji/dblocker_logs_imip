package repository

import (
	"dblocker_logs_server/internal/models"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	DB *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{DB: db}
}

func (r *DeviceRepository) Create(device *models.Device) error {
	return r.DB.Create(device).Error
}

func (r *DeviceRepository) FindAll() ([]models.Device, error) {
	var devices []models.Device
	err := r.DB.Find(&devices).Error
	return devices, err
}

func (r *DeviceRepository) FindByID(id uint) (*models.Device, error) {
	var device models.Device
	err := r.DB.First(&device, id).Error
	return &device, err
}

func (r *DeviceRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Device{}, id).Error
}

func (r *DeviceRepository) Update(device *models.Device) error {
	return r.DB.Save(device).Error
}
