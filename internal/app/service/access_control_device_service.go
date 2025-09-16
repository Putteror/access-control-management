package service

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlDeviceService defines the interface for access control device business logic.
type AccessControlDeviceService interface {
	GetAll(searchQuery schema.AccessControlDeviceSearchQuery) ([]model.AccessControlDevice, error)
	GetByID(id string) (*model.AccessControlDevice, error)
	Save(id string, deviceModel *model.AccessControlDevice) error
	Delete(id string) error
	IsExistName(name string) (bool, error)
	ConvertToResponse(deviceModel *model.AccessControlDevice) (*schema.AccessControlDeviceResponse, error)
}

type accessControlDeviceServiceImpl struct {
	accessControlDeviceRepo repository.AccessControlDeviceRepository
	accessControlServerRepo repository.AccessControlServerRepository
}

// NewAccessControlDeviceService creates a new instance of AccessControlDeviceService.
func NewAccessControlDeviceService(accessControlDeviceRepo repository.AccessControlDeviceRepository, accessControlServerRepo repository.AccessControlServerRepository) AccessControlDeviceService {
	return &accessControlDeviceServiceImpl{
		accessControlDeviceRepo: accessControlDeviceRepo,
		accessControlServerRepo: accessControlServerRepo,
	}
}

// GetAll retrieves all access control devices.
func (s *accessControlDeviceServiceImpl) GetAll(searchQuery schema.AccessControlDeviceSearchQuery) ([]model.AccessControlDevice, error) {
	return s.accessControlDeviceRepo.GetAll(searchQuery)
}

// GetByID retrieves an access control device by its ID.
func (s *accessControlDeviceServiceImpl) GetByID(id string) (*model.AccessControlDevice, error) {
	return s.accessControlDeviceRepo.GetByID(id)
}

// Save creates or updates an access control device.
func (s *accessControlDeviceServiceImpl) Save(id string, deviceModel *model.AccessControlDevice) error {
	// Validate
	if deviceModel.Name == "" {
		return fmt.Errorf("device name cannot be empty")
	}

	if id == "" {
		// Check if name already exists for new device creation
		isExist, err := s.accessControlDeviceRepo.IsExistName(deviceModel.Name)
		if err != nil {
			return fmt.Errorf("failed to check device name existence: %w", err)
		}
		if isExist {
			return fmt.Errorf("device with name '%s' already exists", deviceModel.Name)
		}
		// Create new device
		if err := s.accessControlDeviceRepo.Create(deviceModel); err != nil {
			return fmt.Errorf("failed to create device: %w", err)
		}
	} else {
		// Check existing id
		existingDevice, err := s.accessControlDeviceRepo.GetByID(id)
		if err != nil {
			return fmt.Errorf("failed to get existing device: %w", err)
		}
		if existingDevice == nil {
			return fmt.Errorf("device with ID '%s' not found", id)
		}

		// Update existing device
		if err := s.accessControlDeviceRepo.Update(deviceModel); err != nil {
			return fmt.Errorf("failed to update device: %w", err)
		}
	}

	return nil
}

// Delete deletes an access control device by its ID.
func (s *accessControlDeviceServiceImpl) Delete(id string) error {
	return s.accessControlDeviceRepo.Delete(id)
}

// IsExistName checks if an access control device with the given name exists.
func (s *accessControlDeviceServiceImpl) IsExistName(name string) (bool, error) {
	return s.accessControlDeviceRepo.IsExistName(name)
}

func (s *accessControlDeviceServiceImpl) ConvertToResponse(deviceModel *model.AccessControlDevice) (*schema.AccessControlDeviceResponse, error) {

	var accessControlServerResponse *schema.AccessControlServerInfoResponse
	if deviceModel.AccessControlServerID != nil && *deviceModel.AccessControlServerID != "" {
		server, err := s.accessControlServerRepo.GetByID("ef7d47b5-c9f4-45d1-a8cc-b3fedd61d4c8")
		if err != nil {
			// If the rule is not found, we just return nil for the rule response without an error.
			if err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("failed to find access control rule with ID '%s': %w", *deviceModel.AccessControlServerID, err)
			}
		}
		if server != nil {
			accessControlServerResponse = &schema.AccessControlServerInfoResponse{
				ID:          server.ID,
				Name:        server.Name,
				HostAddress: server.HostAddress,
			}
		}
	}

	response := &schema.AccessControlDeviceResponse{
		ID:                  deviceModel.ID,
		Name:                deviceModel.Name,
		Type:                deviceModel.Type,
		HostAddress:         deviceModel.HostAddress,
		Status:              deviceModel.Status,
		Username:            deviceModel.Username,
		Password:            deviceModel.Password,
		AccessToken:         deviceModel.AccessToken,
		ApiToken:            deviceModel.ApiToken,
		RecordScan:          deviceModel.RecordScan,
		RecordAttendance:    deviceModel.RecordAttendance,
		AllowClockIn:        deviceModel.AllowClockIn,
		AllowClockOut:       deviceModel.AllowClockOut,
		AccessControlServer: accessControlServerResponse,
	}

	fmt.Println(response)

	return response, nil
}
