package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlRuleRepository is the interface for access control rule data access.
type AccessControlRuleRepository interface {
	GetAll(searchQuery schema.AccessControlRuleSearchQuery) ([]model.AccessControlRule, error)
	GetByID(id uuid.UUID) (*model.AccessControlRule, error)
	Create(rule *model.AccessControlRule) error
	Update(rule *model.AccessControlRule) error
	Delete(id uuid.UUID) error
	IsExistName(name string, excludeID uuid.UUID) (bool, error)

	// Group relationship methods
	GetGroupsByRuleID(ruleID uuid.UUID) ([]model.AccessControlRuleGroup, error)
	GetGroupIDsByRuleID(ruleID uuid.UUID) ([]string, error)
	CreateRuleGroups(ruleGroups []model.AccessControlRuleGroup) error
	DeleteRuleGroupsByRuleID(ruleID uuid.UUID, tx *gorm.DB) error
}

// accessControlRuleRepositoryImpl is the implementation of AccessControlRuleRepository.
type accessControlRuleRepositoryImpl struct {
	db *gorm.DB
}

// NewAccessControlRuleRepository creates a new instance of AccessControlRuleRepository.
func NewAccessControlRuleRepository(db *gorm.DB) AccessControlRuleRepository {
	return &accessControlRuleRepositoryImpl{db: db}
}

// --- CRUD Operations ---

// GetAll retrieves all access control rules with pagination.
func (r *accessControlRuleRepositoryImpl) GetAll(searchQuery schema.AccessControlRuleSearchQuery) ([]model.AccessControlRule, error) {
	var rules []model.AccessControlRule

	query := r.db.Model(&model.AccessControlRule{})

	if searchQuery.Name != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery.Name+"%")
	}

	var page int = searchQuery.Page
	var limit int = searchQuery.Limit

	// Check if 'All' is requested, if so, ignore pagination limits
	if !searchQuery.All && page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&rules).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated rules: %w", err)
	}

	return rules, nil
}

// GetByID retrieves a rule by its ID.
func (r *accessControlRuleRepositoryImpl) GetByID(id uuid.UUID) (*model.AccessControlRule, error) {
	var rule model.AccessControlRule
	if err := r.db.First(&rule, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// Create creates a new access control rule record.
func (r *accessControlRuleRepositoryImpl) Create(rule *model.AccessControlRule) error {
	return r.db.Create(rule).Error
}

// Update updates an existing access control rule record.
func (r *accessControlRuleRepositoryImpl) Update(rule *model.AccessControlRule) error {
	return r.db.Save(rule).Error
}

// Delete deletes an access control rule by its ID.
func (r *accessControlRuleRepositoryImpl) Delete(id uuid.UUID) error {
	// ใช้ Transaction เพื่อลบข้อมูลทั้งจากตารางหลักและตารางเชื่อมโยง
	err := r.db.Transaction(func(tx *gorm.DB) error {

		// 1. ลบข้อมูลจากตารางเชื่อมโยง (AccessControlRuleGroup)
		if err := r.DeleteRuleGroupsByRuleID(id, tx); err != nil {
			return err
		}

		// 2. ลบข้อมูลจากตารางหลัก (AccessControlRule)
		if err := tx.Unscoped().Where("id = ?", id).Delete(&model.AccessControlRule{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// IsExistName checks if a rule with the given name exists in the database.
func (r *accessControlRuleRepositoryImpl) IsExistName(name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.AccessControlRule{}).Where("name = ? AND deleted_at IS NULL", name)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check rule existence: %w", err)
	}
	return count > 0, nil
}

// --- Group Relationship Methods ---

// GetGroupsByRuleID retrieves all AccessControlRuleGroup records for a rule ID.
func (r *accessControlRuleRepositoryImpl) GetGroupsByRuleID(ruleID uuid.UUID) ([]model.AccessControlRuleGroup, error) {
	var ruleGroups []model.AccessControlRuleGroup
	err := r.db.Where("access_control_rule_id = ?", ruleID).Find(&ruleGroups).Error
	return ruleGroups, err
}

// GetGroupIDsByRuleID retrieves a list of AccessControlGroupID strings for a rule ID.
func (r *accessControlRuleRepositoryImpl) GetGroupIDsByRuleID(ruleID uuid.UUID) ([]string, error) {
	var groupIDs []string
	err := r.db.Model(&model.AccessControlRuleGroup{}).
		Select("access_control_group_id").
		Where("access_control_rule_id = ?", ruleID).
		Find(&groupIDs).Error
	return groupIDs, err
}

// CreateRuleGroups inserts multiple AccessControlRuleGroup records.
func (r *accessControlRuleRepositoryImpl) CreateRuleGroups(ruleGroups []model.AccessControlRuleGroup) error {
	if len(ruleGroups) == 0 {
		return nil
	}
	// ใช้ r.db โดยตรง (สมมติว่า Service Layer จะจัดการ Transaction)
	return r.db.Create(&ruleGroups).Error
}

// DeleteRuleGroupsByRuleID deletes all AccessControlRuleGroup records for a rule ID.
func (r *accessControlRuleRepositoryImpl) DeleteRuleGroupsByRuleID(ruleID uuid.UUID, tx *gorm.DB) error {
	// ใช้ tx ที่ส่งมาเพื่อรองรับการทำ Transaction
	if tx == nil {
		tx = r.db
	}
	return tx.Unscoped().Where("access_control_rule_id = ?", ruleID).
		Delete(&model.AccessControlRuleGroup{}).Error
}
