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
	IsExistName(name string) (bool, error)
}

// accessControlRuleRepositoryImpl is the implementation of AccessControlRuleRepository.
type accessControlRuleRepositoryImpl struct {
	db *gorm.DB
}

// NewAccessControlRuleRepository creates a new instance of AccessControlRuleRepository.
func NewAccessControlRuleRepository(db *gorm.DB) AccessControlRuleRepository {
	return &accessControlRuleRepositoryImpl{db: db}
}

// GetAll retrieves all access control rules with pagination.
func (r *accessControlRuleRepositoryImpl) GetAll(searchQuery schema.AccessControlRuleSearchQuery) ([]model.AccessControlRule, error) {
	var accessControlRules []model.AccessControlRule

	query := r.db.Model(&model.AccessControlRule{})
	if searchQuery.Name != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery.Name+"%")
	}

	var page int = searchQuery.Page
	var limit int = searchQuery.Limit
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&accessControlRules).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated rules: %w", err)
	}
	return accessControlRules, nil
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
	return r.db.Unscoped().Where("id = ?", id).Delete(&model.AccessControlRule{}).Error
}

// IsExistName checks if a rule with the given name exists in the database.
func (r *accessControlRuleRepositoryImpl) IsExistName(name string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.AccessControlRule{}).Where("name = ? AND deleted_at IS NULL", name).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check rule existence: %w", err)
	}
	return count > 0, nil
}
