package service

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
)

type DeviceService interface {
	CreateDevice(device *model.Device) error
	GetAllDevice() ([]model.Device, error)
}

type deviceServiceImpl struct {
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) DeviceService {
	return &deviceServiceImpl{repo: repo}
}

func (s *deviceServiceImpl) CreateDevice(device *model.Device) error {
	// Business Logic
	if device.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return s.repo.Create(device)
}

func (s *deviceServiceImpl) GetAllDevice() ([]model.Device, error) {
	return s.repo.FindAll()
}
