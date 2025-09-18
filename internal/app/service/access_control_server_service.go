package service

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlServerService defines the interface for access control server business logic.
type AccessControlServerService interface {
	GetAll(searchQuery schema.AccessControlServerSearchQuery) ([]model.AccessControlServer, error)
	GetByID(id string) (*model.AccessControlServer, error)
	Create(bodyRequest *schema.AccessControlServerRequest) (*model.AccessControlServer, error)
	Update(id string, bodyRequest *schema.AccessControlServerRequest) (*model.AccessControlServer, error)
	PartialUpdate(id string, bodyRequest *schema.AccessControlServerRequest) (*model.AccessControlServer, error)
	Delete(id string) error
	ConvertToResponse(serverModel *model.AccessControlServer) (*schema.AccessControlServerResponse, error)
}

type accessControlServerServiceImpl struct {
	accessControlServerRepo repository.AccessControlServerRepository
}

// NewAccessControlServerService creates a new instance of AccessControlServerService.
func NewAccessControlServerService(accessControlServerRepo repository.AccessControlServerRepository) AccessControlServerService {
	return &accessControlServerServiceImpl{
		accessControlServerRepo: accessControlServerRepo,
	}
}

// GetAll retrieves all access control servers.
func (s *accessControlServerServiceImpl) GetAll(searchQuery schema.AccessControlServerSearchQuery) ([]model.AccessControlServer, error) {
	return s.accessControlServerRepo.GetAll(searchQuery)
}

// GetByID retrieves an access control server by its ID.
func (s *accessControlServerServiceImpl) GetByID(id string) (*model.AccessControlServer, error) {
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.accessControlServerRepo.GetByID(id_uuid)
}

// Create creates a new access control server.
func (s *accessControlServerServiceImpl) Create(bodyRequest *schema.AccessControlServerRequest) (*model.AccessControlServer, error) {

	// Set default value
	bodyRequest, err := s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	// Create Model
	serverModel := &model.AccessControlServer{
		Name:        *bodyRequest.Name,
		Type:        *bodyRequest.Type,
		HostAddress: *bodyRequest.HostAddress,
		Username:    bodyRequest.Username,
		Password:    bodyRequest.Password,
		AccessToken: bodyRequest.AccessToken,
		ApiToken:    bodyRequest.ApiToken,
		Status:      *bodyRequest.Status,
	}
	// Validate
	validateDuplicateErr := s.validateBodyRequest(*bodyRequest, nil)
	if validateDuplicateErr != nil {
		return nil, validateDuplicateErr
	}
	// Create server
	if err := s.accessControlServerRepo.Create(serverModel); err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	return serverModel, nil
}

// Update updates an existing access control server.
func (s *accessControlServerServiceImpl) Update(id string, bodyRequest *schema.AccessControlServerRequest) (*model.AccessControlServer, error) {

	// Check have item
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	serverModel, err := s.accessControlServerRepo.GetByID(id_uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing server: %w", err)
	}
	if serverModel == nil {
		return nil, fmt.Errorf("server with ID '%s' not found", id)
	}

	// Validate with old model
	bodyRequest, err = s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	err = s.validateBodyRequest(*bodyRequest, serverModel)
	if err != nil {
		return nil, err
	}
	// Update model
	serverModel.Name = *bodyRequest.Name
	serverModel.Type = *bodyRequest.Type
	serverModel.HostAddress = *bodyRequest.HostAddress
	serverModel.Username = bodyRequest.Username
	serverModel.Password = bodyRequest.Password
	serverModel.AccessToken = bodyRequest.AccessToken
	serverModel.ApiToken = bodyRequest.ApiToken
	serverModel.Status = *bodyRequest.Status
	// Update existing server
	if err := s.accessControlServerRepo.Update(serverModel); err != nil {
		return nil, fmt.Errorf("failed to update server: %w", err)
	}

	return serverModel, nil
}

// PartialUpdate performs a partial update on an existing access control server.
func (s *accessControlServerServiceImpl) PartialUpdate(id string, bodyRequest *schema.AccessControlServerRequest) (*model.AccessControlServer, error) {

	// Get model from id
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	serverModel, err := s.accessControlServerRepo.GetByID(id_uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing server: %w", err)
	}
	if serverModel == nil {
		return nil, fmt.Errorf("server with ID '%s' not found", id)
	}

	// Validate with old model
	validateDuplicateErr := s.validateBodyRequest(*bodyRequest, serverModel)
	if validateDuplicateErr != nil {
		return nil, validateDuplicateErr
	}
	// Update model
	if bodyRequest.Name != nil {
		serverModel.Name = *bodyRequest.Name
	}
	if bodyRequest.Type != nil {
		serverModel.Type = *bodyRequest.Type
	}
	if bodyRequest.HostAddress != nil {
		serverModel.HostAddress = *bodyRequest.HostAddress
	}
	if bodyRequest.Username != nil {
		serverModel.Username = bodyRequest.Username
	}
	if bodyRequest.Password != nil {
		serverModel.Password = bodyRequest.Password
	}
	if bodyRequest.AccessToken != nil {
		serverModel.AccessToken = bodyRequest.AccessToken
	}
	if bodyRequest.ApiToken != nil {
		serverModel.ApiToken = bodyRequest.ApiToken
	}
	if bodyRequest.Status != nil {
		serverModel.Status = *bodyRequest.Status
	}
	// Update existing server
	if err := s.accessControlServerRepo.Update(serverModel); err != nil {
		return nil, fmt.Errorf("failed to update server: %w", err)
	}

	return serverModel, nil
}

// Delete deletes an access control server by its ID.
func (s *accessControlServerServiceImpl) Delete(id string) error {
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}
	_, err = s.accessControlServerRepo.GetByID(id_uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("server with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get server by ID: %w", err)
	}
	return s.accessControlServerRepo.Delete(id_uuid)
}

// ConvertToResponse converts a server model to a response schema.
func (s *accessControlServerServiceImpl) ConvertToResponse(serverModel *model.AccessControlServer) (*schema.AccessControlServerResponse, error) {

	response := &schema.AccessControlServerResponse{
		ID:          serverModel.ID.String(),
		Name:        serverModel.Name,
		Type:        serverModel.Type,
		HostAddress: serverModel.HostAddress,
		Status:      serverModel.Status,
		Username:    serverModel.Username,
		AccessToken: serverModel.AccessToken,
		ApiToken:    serverModel.ApiToken,
	}
	// Omit password from the response for security reasons
	if serverModel.Password != nil && strings.TrimSpace(*serverModel.Password) != "" {
		emptyPassword := ""
		response.Password = &emptyPassword
	}

	return response, nil
}

// validateAndSetDefaultValues validates request data and sets default values.
func (s *accessControlServerServiceImpl) validateAndSetDefaultValues(bodyRequest *schema.AccessControlServerRequest) (*schema.AccessControlServerRequest, error) {

	if bodyRequest.Name == nil || *bodyRequest.Name == "" {
		return nil, fmt.Errorf("server name cannot be empty")
	}
	if bodyRequest.HostAddress == nil || *bodyRequest.HostAddress == "" {
		return nil, fmt.Errorf("server host address cannot be empty")
	}
	if bodyRequest.Type == nil || *bodyRequest.Type == "" {
		return nil, fmt.Errorf("server type cannot be empty")
	}

	if bodyRequest.Status == nil {
		bodyRequest.Status = new(string)
		*bodyRequest.Status = "active"
	}
	return bodyRequest, nil
}

// validateBodyRequest checks for duplicate name or host address.
func (s *accessControlServerServiceImpl) validateBodyRequest(bodyRequest schema.AccessControlServerRequest, serverModel *model.AccessControlServer) error {

	excludeID := uuid.Nil
	if serverModel != nil {
		excludeID = serverModel.ID
	}

	// Check duplicate name
	if bodyRequest.Name != nil {
		isExistName, err := s.accessControlServerRepo.IsExistName(*bodyRequest.Name, excludeID)
		if err != nil {
			return fmt.Errorf("failed to check server name existence: %w", err)
		}
		if isExistName {
			return fmt.Errorf("server with name '%s' already exists", *bodyRequest.Name)
		}
	}
	// Check duplicate host address
	if bodyRequest.HostAddress != nil {
		isExistHostAddress, err := s.accessControlServerRepo.IsExistHostAddress(*bodyRequest.HostAddress, excludeID)
		if err != nil {
			return fmt.Errorf("failed to check server host address existence: %w", err)
		}
		if isExistHostAddress {
			return fmt.Errorf("server with host address '%s' already exists", *bodyRequest.HostAddress)
		}
	}

	return nil
}
