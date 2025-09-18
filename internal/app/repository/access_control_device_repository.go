package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlDeviceRepository is the interface for access control device data access.
type AccessControlDeviceRepository interface {
	GetAll(searchQuery schema.AccessControlDeviceSearchQuery) ([]model.AccessControlDevice, error)
	GetByID(id uuid.UUID) (*model.AccessControlDevice, error)
	Create(device *model.AccessControlDevice) error
	Update(device *model.AccessControlDevice) error
	Delete(id uuid.UUID) error
	IsExistName(name string, excludeID uuid.UUID) (bool, error)
	IsExistHostAddress(hostAddress string, excludeID uuid.UUID) (bool, error)
}

// accessControlDeviceRepositoryImpl is the implementation of AccessControlDeviceRepository.
type accessControlDeviceRepositoryImpl struct {
	db *gorm.DB
}

// NewAccessControlDeviceRepository creates a new instance of AccessControlDeviceRepository.
func NewAccessControlDeviceRepository(db *gorm.DB) AccessControlDeviceRepository {
	return &accessControlDeviceRepositoryImpl{db: db}
}

// GetAll retrieves all access control devices with pagination.
func (r *accessControlDeviceRepositoryImpl) GetAll(searchQuery schema.AccessControlDeviceSearchQuery) ([]model.AccessControlDevice, error) {
	var devices []model.AccessControlDevice

	query := r.db.Model(&model.AccessControlDevice{})

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
	if err := query.Offset(offset).Limit(limit).Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated devices: %w", err)
	}

	return devices, nil
}

// GetByID retrieves a device by its ID.
func (r *accessControlDeviceRepositoryImpl) GetByID(id uuid.UUID) (*model.AccessControlDevice, error) {
	var device model.AccessControlDevice
	if err := r.db.First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// Create creates a new access control device record.
func (r *accessControlDeviceRepositoryImpl) Create(device *model.AccessControlDevice) error {
	return r.db.Create(device).Error
}

// Update updates an existing access control device record.
func (r *accessControlDeviceRepositoryImpl) Update(device *model.AccessControlDevice) error {
	return r.db.Save(device).Error
}

// Delete deletes an access control device by its ID.
func (r *accessControlDeviceRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Unscoped().Where("id = ?", id).Delete(&model.AccessControlDevice{}).Error
}

// IsExistName checks if a device with the given name exists in the database.
func (r *accessControlDeviceRepositoryImpl) IsExistName(name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.AccessControlDevice{}).Where("name = ? AND deleted_at IS NULL", name)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check device existence: %w", err)
	}
	return count > 0, nil
}

// IsExistHostAddress
func (r *accessControlDeviceRepositoryImpl) IsExistHostAddress(hostAddress string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.AccessControlDevice{}).Where("host_address = ? AND deleted_at IS NULL", hostAddress)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check device existence: %w", err)
	}
	return count > 0, nil
}
