package repository

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlServerRepository is the interface for access control server data access.
type AccessControlServerRepository interface {
	GetAll(searchQuery schema.AccessControlServerSearchQuery) ([]model.AccessControlServer, error)
	GetByID(id string) (*model.AccessControlServer, error)
	Create(server *model.AccessControlServer) error
	Update(server *model.AccessControlServer) error
	Delete(id string) error
	IsExistName(name string) (bool, error)
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

	var page int = searchQuery.Page
	var limit int = searchQuery.Limit
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&servers).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated servers: %w", err)
	}
	return servers, nil
}

// GetByID retrieves a server by its ID.
func (r *accessControlServerRepositoryImpl) GetByID(id string) (*model.AccessControlServer, error) {
	var server model.AccessControlServer
	println("repo")
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
func (r *accessControlServerRepositoryImpl) Delete(id string) error {
	return r.db.Unscoped().Where("id = ?", id).Delete(&model.AccessControlServer{}).Error
}

// IsExistName checks if a server with the given name exists in the database.
func (r *accessControlServerRepositoryImpl) IsExistName(name string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.AccessControlServer{}).Where("name = ? AND deleted_at IS NULL", name).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check server existence: %w", err)
	}
	return count > 0, nil
}
