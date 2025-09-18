package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlDeviceService defines the interface for access control device business logic.
type AccessControlDeviceService interface {
	GetAll(searchQuery schema.AccessControlDeviceSearchQuery) ([]model.AccessControlDevice, error)
	GetByID(id string) (*model.AccessControlDevice, error)
	Create(bodyRequest *schema.AccessControlDeviceRequest) (*model.AccessControlDevice, error)
	Update(id string, bodyRequest *schema.AccessControlDeviceRequest) (*model.AccessControlDevice, error)
	PartialUpdate(id string, bodyRequest *schema.AccessControlDeviceRequest) (*model.AccessControlDevice, error)
	Delete(id string) error
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
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.accessControlDeviceRepo.GetByID(id_uuid)
}

// Save creates or updates an access control device.
func (s *accessControlDeviceServiceImpl) Create(bodyRequest *schema.AccessControlDeviceRequest) (*model.AccessControlDevice, error) {

	// Validate and Set default value
	bodyRequest, err := s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	validateDuplicateErr := s.validateBodyRequest(*bodyRequest, nil)
	if validateDuplicateErr != nil {
		return nil, validateDuplicateErr
	}
	// Create Model
	deviceModel := &model.AccessControlDevice{
		Name:                  *bodyRequest.Name,
		Type:                  *bodyRequest.Type,
		HostAddress:           *bodyRequest.HostAddress,
		Username:              bodyRequest.Username,
		Password:              bodyRequest.Password,
		AccessToken:           bodyRequest.AccessToken,
		ApiToken:              bodyRequest.ApiToken,
		RecordScan:            *bodyRequest.RecordScan,
		RecordAttendance:      *bodyRequest.RecordAttendance,
		AllowClockIn:          *bodyRequest.AllowClockIn,
		AllowClockOut:         *bodyRequest.AllowClockOut,
		Status:                *bodyRequest.Status,
		AccessControlServerID: bodyRequest.AccessControlServerID,
	}

	// Update existing device
	if err := s.accessControlDeviceRepo.Create(deviceModel); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return deviceModel, nil
}

func (s *accessControlDeviceServiceImpl) Update(id string, bodyRequest *schema.AccessControlDeviceRequest) (*model.AccessControlDevice, error) {

	// Check have item
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	deviceModel, err := s.accessControlDeviceRepo.GetByID(id_uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing device: %w", err)
	}
	if deviceModel == nil {
		return nil, fmt.Errorf("device with ID '%s' not found", id)
	}

	// Validate with old model
	bodyRequest, err = s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	err = s.validateBodyRequest(*bodyRequest, deviceModel)
	if err != nil {
		return nil, err
	}
	// Update model
	deviceModel.Name = *bodyRequest.Name
	deviceModel.Type = *bodyRequest.Type
	deviceModel.HostAddress = *bodyRequest.HostAddress
	deviceModel.Username = bodyRequest.Username
	deviceModel.Password = bodyRequest.Password
	deviceModel.AccessToken = bodyRequest.AccessToken
	deviceModel.ApiToken = bodyRequest.ApiToken
	deviceModel.RecordScan = *bodyRequest.RecordScan
	deviceModel.RecordAttendance = *bodyRequest.RecordAttendance
	deviceModel.AllowClockIn = *bodyRequest.AllowClockIn
	deviceModel.AllowClockOut = *bodyRequest.AllowClockOut
	deviceModel.Status = *bodyRequest.Status
	deviceModel.AccessControlServerID = bodyRequest.AccessControlServerID
	// Update existing device
	if err := s.accessControlDeviceRepo.Update(deviceModel); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return deviceModel, nil
}

func (s *accessControlDeviceServiceImpl) PartialUpdate(id string, bodyRequest *schema.AccessControlDeviceRequest) (*model.AccessControlDevice, error) {

	// Get model from id
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	deviceModel, err := s.accessControlDeviceRepo.GetByID(id_uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing device: %w", err)
	}
	if deviceModel == nil {
		return nil, fmt.Errorf("device with ID '%s' not found", id)
	}
	// Validate with old model
	validateDuplicateErr := s.validateBodyRequest(*bodyRequest, deviceModel)
	if validateDuplicateErr != nil {
		return nil, validateDuplicateErr
	}
	// Update model
	if bodyRequest.Name != nil {
		deviceModel.Name = *bodyRequest.Name
	}
	if bodyRequest.Type != nil {
		deviceModel.Type = *bodyRequest.Type
	}
	if bodyRequest.HostAddress != nil {
		deviceModel.HostAddress = *bodyRequest.HostAddress
	}
	if bodyRequest.Username != nil {
		deviceModel.Username = bodyRequest.Username
	}
	if bodyRequest.Password != nil {
		deviceModel.Password = bodyRequest.Password
	}
	if bodyRequest.AccessToken != nil {
		deviceModel.AccessToken = bodyRequest.AccessToken
	}
	if bodyRequest.ApiToken != nil {
		deviceModel.ApiToken = bodyRequest.ApiToken
	}
	if bodyRequest.RecordScan != nil {
		deviceModel.RecordScan = *bodyRequest.RecordScan
	}
	if bodyRequest.RecordAttendance != nil {
		deviceModel.RecordAttendance = *bodyRequest.RecordAttendance
	}
	if bodyRequest.AllowClockIn != nil {
		deviceModel.AllowClockIn = *bodyRequest.AllowClockIn
	}
	if bodyRequest.AllowClockOut != nil {
		deviceModel.AllowClockOut = *bodyRequest.AllowClockOut
	}
	if bodyRequest.Status != nil {
		deviceModel.Status = *bodyRequest.Status
	}
	if bodyRequest.AccessControlServerID != nil {
		deviceModel.AccessControlServerID = bodyRequest.AccessControlServerID
	}
	// Update existing device
	if err := s.accessControlDeviceRepo.Update(deviceModel); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return deviceModel, nil
}

// Delete deletes an access control device by its ID.
func (s *accessControlDeviceServiceImpl) Delete(id string) error {
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}
	_, err = s.accessControlDeviceRepo.GetByID(id_uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("device with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get device by ID: %w", err)
	}

	return s.accessControlDeviceRepo.Delete(id_uuid)
}

func (s *accessControlDeviceServiceImpl) ConvertToResponse(deviceModel *model.AccessControlDevice) (*schema.AccessControlDeviceResponse, error) {

	var accessControlServerResponse *schema.AccessControlServerInfoResponse
	if deviceModel.AccessControlServerID != nil && *deviceModel.AccessControlServerID != "" {
		access_control_server_uuid, err := uuid.Parse(*deviceModel.AccessControlServerID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse access control server ID: %w", err)
		}
		server, err := s.accessControlServerRepo.GetByID(access_control_server_uuid)
		if err != nil {
			// If the rule is not found, we just return nil for the rule response without an error.
			if err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("failed to find access control rule with ID '%s': %w", *deviceModel.AccessControlServerID, err)
			}
		}
		if server != nil {
			accessControlServerResponse = &schema.AccessControlServerInfoResponse{
				ID:          server.ID.String(),
				Name:        server.Name,
				HostAddress: server.HostAddress,
			}
		}
	}

	response := &schema.AccessControlDeviceResponse{
		ID:                  deviceModel.ID.String(),
		Name:                deviceModel.Name,
		Type:                deviceModel.Type,
		HostAddress:         deviceModel.HostAddress,
		Status:              deviceModel.Status,
		Username:            deviceModel.Username,
		RecordScan:          deviceModel.RecordScan,
		RecordAttendance:    deviceModel.RecordAttendance,
		AllowClockIn:        deviceModel.AllowClockIn,
		AllowClockOut:       deviceModel.AllowClockOut,
		AccessControlServer: accessControlServerResponse,
	}

	return response, nil
}

func (s *accessControlDeviceServiceImpl) validateAndSetDefaultValues(bodyRequest *schema.AccessControlDeviceRequest) (*schema.AccessControlDeviceRequest, error) {

	if bodyRequest.Name == nil || *bodyRequest.Name == "" {
		return nil, fmt.Errorf("device name cannot be empty")
	}
	if bodyRequest.HostAddress == nil || *bodyRequest.HostAddress == "" {
		return nil, fmt.Errorf("device host address cannot be empty")
	}
	if bodyRequest.Type == nil || *bodyRequest.Type == "" {
		return nil, fmt.Errorf("device type cannot be empty")
	}
	if bodyRequest.RecordScan == nil {
		bodyRequest.RecordScan = new(bool)
		*bodyRequest.RecordScan = false
	}
	if bodyRequest.RecordAttendance == nil {
		bodyRequest.RecordAttendance = new(bool)
		*bodyRequest.RecordAttendance = false
	}
	if bodyRequest.AllowClockIn == nil {
		bodyRequest.AllowClockIn = new(bool)
		*bodyRequest.AllowClockIn = false
	}
	if bodyRequest.AllowClockOut == nil {
		bodyRequest.AllowClockOut = new(bool)
		*bodyRequest.AllowClockOut = false
	}
	if bodyRequest.Status == nil {
		bodyRequest.Status = new(string)
		*bodyRequest.Status = "active"
	}
	return bodyRequest, nil
}

func (s *accessControlDeviceServiceImpl) validateBodyRequest(bodyRequest schema.AccessControlDeviceRequest, deviceModel *model.AccessControlDevice) error {

	excludeID := uuid.Nil
	if deviceModel != nil {
		excludeID = deviceModel.ID
	}

	// Checl ExistId
	if bodyRequest.AccessControlServerID != nil {
		server_uuid, err := uuid.Parse(*bodyRequest.AccessControlServerID)
		if err != nil {
			return fmt.Errorf("invalid access control server ID")
		}
		_, err = s.accessControlServerRepo.GetByID(server_uuid)
		if err != nil {
			return fmt.Errorf("failed to get access control server by ID: %w", err)
		}
	}

	// Check duplicate
	if bodyRequest.Name != nil {
		isExistName, err := s.accessControlDeviceRepo.IsExistName(*bodyRequest.Name, excludeID)
		if err != nil {
			return fmt.Errorf("failed to check device name existence: %w", err)
		}
		if isExistName {
			return fmt.Errorf("device with name '%s' already exists", *bodyRequest.Name)
		}
	}
	if bodyRequest.HostAddress != nil {
		isExistHostAddress, err := s.accessControlDeviceRepo.IsExistHostAddress(*bodyRequest.HostAddress, excludeID)
		if err != nil {
			return fmt.Errorf("failed to check device name existence: %w", err)
		}
		if isExistHostAddress {
			return fmt.Errorf("device with name '%s' already exists", *bodyRequest.HostAddress)
		}
	}

	return nil
}
