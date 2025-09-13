package repository

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"gorm.io/gorm"
)

// AccessControlRuleRepository is the interface for rule data access.
type AccessControlRuleRepository interface {
	// IsExist checks if a rule with the given ID exists in the database.
	IsExistID(id string) (bool, error)
}

// accessControlRuleRepositoryImpl is the implementation of AccessControlRuleRepository.
type accessControlRuleRepositoryImpl struct {
	db *gorm.DB
}

// NewAccessControlRuleRepository creates a new instance of AccessControlRuleRepository.
func NewAccessControlRuleRepository(db *gorm.DB) AccessControlRuleRepository {
	return &accessControlRuleRepositoryImpl{db: db}
}

// IsExist checks if a rule with the given ID exists in the database.
func (r *accessControlRuleRepositoryImpl) IsExistID(id string) (bool, error) {
	var count int64
	result := r.db.Model(&model.AccessControlRule{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check rule existence: %w", result.Error)
	}

	return count > 0, nil
}
