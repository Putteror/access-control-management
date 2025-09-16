package service

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
)

// AccessControlServerService defines the interface for access control server business logic.
type AccessControlServerService interface {
	GetAll(searchQuery schema.AccessControlServerSearchQuery) ([]model.AccessControlServer, error)
	GetByID(id string) (*model.AccessControlServer, error)
	Save(id string, serverModel *model.AccessControlServer) error
	Delete(id string) error
	IsExistName(name string) (bool, error)
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
	return s.accessControlServerRepo.GetByID(id)
}

// Save creates or updates an access control server.
func (s *accessControlServerServiceImpl) Save(id string, serverModel *model.AccessControlServer) error {
	// Validate
	if serverModel.Name == "" {
		return fmt.Errorf("server name cannot be empty")
	}

	if id == "" {
		// Check if name already exists for new server creation
		isExist, err := s.accessControlServerRepo.IsExistName(serverModel.Name)
		if err != nil {
			return fmt.Errorf("failed to check server name existence: %w", err)
		}
		if isExist {
			return fmt.Errorf("server with name '%s' already exists", serverModel.Name)
		}
		// Create new server
		if err := s.accessControlServerRepo.Create(serverModel); err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}
	} else {
		// Check existing id
		existingServer, err := s.accessControlServerRepo.GetByID(id)
		if err != nil {
			return fmt.Errorf("failed to get existing server: %w", err)
		}
		if existingServer == nil {
			return fmt.Errorf("server with ID '%s' not found", id)
		}

		// Update existing server
		if err := s.accessControlServerRepo.Update(serverModel); err != nil {
			return fmt.Errorf("failed to update server: %w", err)
		}
	}

	return nil
}

// Delete deletes an access control server by its ID.
func (s *accessControlServerServiceImpl) Delete(id string) error {
	return s.accessControlServerRepo.Delete(id)
}

// IsExistName checks if an access control server with the given name exists.
func (s *accessControlServerServiceImpl) IsExistName(name string) (bool, error) {
	return s.accessControlServerRepo.IsExistName(name)
}
