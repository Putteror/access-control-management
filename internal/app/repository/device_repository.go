package repository

import (
	"github.com/putteror/access-control-management/internal/app/model"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	Create(device *model.Device) error
	FindAll() ([]model.Device, error)
}

type deviceRepositoryImpl struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepositoryImpl{db: db}
}

func (r *deviceRepositoryImpl) Create(device *model.Device) error {
	return r.db.Create(device).Error
}

func (r *deviceRepositoryImpl) FindAll() ([]model.Device, error) {
	var devices []model.Device
	if err := r.db.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
