package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlGroupRepository is the interface for access control group data access.
type AccessControlGroupRepository interface {
	GetAll(searchQuery schema.AccessControlGroupSearchQuery) ([]model.AccessControlGroup, error)
	GetByID(id uuid.UUID) (*model.AccessControlGroup, error)
	Create(group *model.AccessControlGroup) error
	Update(group *model.AccessControlGroup) error
	Delete(id uuid.UUID) error
	IsExistName(name string, excludeID uuid.UUID) (bool, error)

	// Device relationship methods
	GetDevicesByGroupID(groupID uuid.UUID) ([]model.AccessControlGroupDevice, error)
	GetDeviceIDsByGroupID(groupID uuid.UUID) ([]string, error)
	CreateGroupDevices(groupDevices []model.AccessControlGroupDevice) error
	DeleteGroupDevicesByGroupID(groupID uuid.UUID, tx *gorm.DB) error
	GetAccessControlGroupScheduleByGroupID(groupID string) ([]model.AccessControlGroupSchedule, error)
	CreateAccessControlGroupSchedule(groupSchedules []model.AccessControlGroupSchedule) error
	DeleteAccessControlGroupScheduleByGroupID(groupID uuid.UUID, tx *gorm.DB) error
}

// accessControlGroupRepositoryImpl is the implementation of AccessControlGroupRepository.
type accessControlGroupRepositoryImpl struct {
	db *gorm.DB
}

// NewAccessControlGroupRepository creates a new instance of AccessControlGroupRepository.
func NewAccessControlGroupRepository(db *gorm.DB) AccessControlGroupRepository {
	return &accessControlGroupRepositoryImpl{db: db}
}

// GetAll retrieves all access control groups with pagination.
func (r *accessControlGroupRepositoryImpl) GetAll(searchQuery schema.AccessControlGroupSearchQuery) ([]model.AccessControlGroup, error) {
	var groups []model.AccessControlGroup

	query := r.db.Model(&model.AccessControlGroup{})

	if searchQuery.Name != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery.Name+"%")
	}

	var page int = searchQuery.Page
	var limit int = searchQuery.Limit
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated groups: %w", err)
	}

	return groups, nil
}

// GetByID retrieves a group by its ID.
func (r *accessControlGroupRepositoryImpl) GetByID(id uuid.UUID) (*model.AccessControlGroup, error) {
	var group model.AccessControlGroup
	if err := r.db.First(&group, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// Create creates a new access control group record.
func (r *accessControlGroupRepositoryImpl) Create(group *model.AccessControlGroup) error {
	return r.db.Create(group).Error
}

// Update updates an existing access control group record.
func (r *accessControlGroupRepositoryImpl) Update(group *model.AccessControlGroup) error {
	return r.db.Save(group).Error
}

// Delete deletes an access control group by its ID.
func (r *accessControlGroupRepositoryImpl) Delete(id uuid.UUID) error {
	// ใช้ Transaction เพื่อลบข้อมูลทั้งจากตารางหลักและตารางเชื่อมโยง
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. ลบข้อมูลจากตารางเชื่อมโยง (AccessControlGroupDevice)
		if err := r.DeleteGroupDevicesByGroupID(id, tx); err != nil {
			return err
		}
		// 2. ลบข้อมูลจากตารางหลัก (AccessControlGroup)
		if err := tx.Unscoped().Where("id = ?", id).Delete(&model.AccessControlGroup{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// IsExistName checks if a group with the given name exists in the database.
func (r *accessControlGroupRepositoryImpl) IsExistName(name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.AccessControlGroup{}).Where("name = ? AND deleted_at IS NULL", name)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check group existence: %w", err)
	}
	return count > 0, nil
}

// --- Device Relationship Methods ---

// GetDevicesByGroupID retrieves all AccessControlGroupDevice records for a group ID.
func (r *accessControlGroupRepositoryImpl) GetDevicesByGroupID(groupID uuid.UUID) ([]model.AccessControlGroupDevice, error) {
	var groupDevices []model.AccessControlGroupDevice
	err := r.db.Where("access_control_group_id = ?", groupID).Find(&groupDevices).Error
	return groupDevices, err
}

// GetDeviceIDsByGroupID retrieves a list of AccessControlDeviceID strings for a group ID.
func (r *accessControlGroupRepositoryImpl) GetDeviceIDsByGroupID(groupID uuid.UUID) ([]string, error) {
	var deviceIDs []string
	err := r.db.Model(&model.AccessControlGroupDevice{}).
		Select("access_control_device_id").
		Where("access_control_group_id = ?", groupID).
		Find(&deviceIDs).Error
	return deviceIDs, err
}

// CreateGroupDevices inserts multiple AccessControlGroupDevice records.
func (r *accessControlGroupRepositoryImpl) CreateGroupDevices(groupDevices []model.AccessControlGroupDevice) error {
	if len(groupDevices) == 0 {
		return nil
	}
	return r.db.Create(&groupDevices).Error
}

// DeleteGroupDevicesByGroupID deletes all AccessControlGroupDevice records for a group ID.
func (r *accessControlGroupRepositoryImpl) DeleteGroupDevicesByGroupID(groupID uuid.UUID, tx *gorm.DB) error {
	// ใช้ tx ที่ส่งมาเพื่อรองรับการทำ Transaction
	if tx == nil {
		tx = r.db
	}
	return tx.Unscoped().Where("access_control_group_id = ?", groupID).
		Delete(&model.AccessControlGroupDevice{}).Error
}

// GetAccessControlGroupScheduleByGroupID
func (r *accessControlGroupRepositoryImpl) GetAccessControlGroupScheduleByGroupID(groupID string) ([]model.AccessControlGroupSchedule, error) {
	var groupSchedules []model.AccessControlGroupSchedule
	err := r.db.Where("access_control_group_id = ?", groupID).Find(&groupSchedules).Error
	return groupSchedules, err
}

// CreateGroupSchedule
func (r *accessControlGroupRepositoryImpl) CreateAccessControlGroupSchedule(groupSchedules []model.AccessControlGroupSchedule) error {
	if len(groupSchedules) == 0 {
		return nil
	}
	return r.db.Create(&groupSchedules).Error
}

// DeleteGroupSchecule
func (r *accessControlGroupRepositoryImpl) DeleteAccessControlGroupScheduleByGroupID(groupID uuid.UUID, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Unscoped().Where("access_control_group_id = ?", groupID).
		Delete(&model.AccessControlGroupSchedule{}).Error
}
