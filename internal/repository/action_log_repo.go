package repository

import (
	"dblocker_logs_server/internal/models"

	"gorm.io/gorm"
)

type ActionLogRepository struct {
	DB *gorm.DB
}

func NewActionLogRepository(db *gorm.DB) *ActionLogRepository {
	return &ActionLogRepository{DB: db}
}

func (r *ActionLogRepository) Create(log *models.ActionLog) error {
	return r.DB.Create(log).Error
}

func (r *ActionLogRepository) FindAll() ([]models.ActionLog, error) {
	var logs []models.ActionLog
	err := r.DB.Find(&logs).Error
	return logs, err
}

func (r *ActionLogRepository) FindByID(id uint) (*models.ActionLog, error) {
	var log models.ActionLog
	err := r.DB.First(&log, id).Error
	return &log, err
}

func (r *ActionLogRepository) Update(log *models.ActionLog) error {
	return r.DB.Save(log).Error
}

func (r *ActionLogRepository) Delete(id uint) error {
	return r.DB.Delete(&models.ActionLog{}, id).Error
}
