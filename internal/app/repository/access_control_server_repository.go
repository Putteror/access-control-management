package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlServerRepository is the interface for access control server data access.
type AccessControlServerRepository interface {
	GetAll(searchQuery schema.AccessControlServerSearchQuery) ([]model.AccessControlServer, error)
	GetByID(id uuid.UUID) (*model.AccessControlServer, error)
	Create(server *model.AccessControlServer) error
	Update(server *model.AccessControlServer) error
	Delete(id uuid.UUID) error
	IsExistName(name string, excludeID uuid.UUID) (bool, error)
	IsExistHostAddress(hostAddress string, excludeID uuid.UUID) (bool, error)
}

// accessControlServerRepositoryImpl is the implementation of AccessControlServerRepository.
type accessControlServerRepositoryImpl struct {
	db *gorm.DB
}

// NewAccessControlServerRepository creates a new instance of AccessControlServerRepository.
func NewAccessControlServerRepository(db *gorm.DB) AccessControlServerRepository {
	return &accessControlServerRepositoryImpl{db: db}
}

// GetAll retrieves all access control servers with pagination.
func (r *accessControlServerRepositoryImpl) GetAll(searchQuery schema.AccessControlServerSearchQuery) ([]model.AccessControlServer, error) {
	var servers []model.AccessControlServer

	query := r.db.Model(&model.AccessControlServer{})

	if searchQuery.Name != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery.Name+"%")
	}

	if searchQuery.Type != "" {
		query = query.Where("type ILIKE ?", "%"+searchQuery.Type+"%")
	}

	if searchQuery.HostAddress != "" {
		query = query.Where("host_address ILIKE ?", "%"+searchQuery.HostAddress+"%")
	}

	var page int = searchQuery.Page
	var limit int = searchQuery.Limit
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&servers).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated servers: %w", err)
	}

	return servers, nil
}

// GetByID retrieves a server by its ID.
func (r *accessControlServerRepositoryImpl) GetByID(id uuid.UUID) (*model.AccessControlServer, error) {
	var server model.AccessControlServer
	if err := r.db.First(&server, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

// Create creates a new access control server record.
func (r *accessControlServerRepositoryImpl) Create(server *model.AccessControlServer) error {
	return r.db.Create(server).Error
}

// Update updates an existing access control server record.
func (r *accessControlServerRepositoryImpl) Update(server *model.AccessControlServer) error {
	return r.db.Save(server).Error
}

// Delete deletes an access control server by its ID.
func (r *accessControlServerRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Unscoped().Where("id = ?", id).Delete(&model.AccessControlServer{}).Error
}

// IsExistName checks if a server with the given name exists in the database.
func (r *accessControlServerRepositoryImpl) IsExistName(name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.AccessControlServer{}).Where("name = ? AND deleted_at IS NULL", name)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check server existence: %w", err)
	}
	return count > 0, nil
}

// IsExistHostAddress checks if a server with the given host address exists in the database.
func (r *accessControlServerRepositoryImpl) IsExistHostAddress(hostAddress string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.AccessControlServer{}).Where("host_address = ? AND deleted_at IS NULL", hostAddress)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check server existence: %w", err)
	}
	return count > 0, nil
}
