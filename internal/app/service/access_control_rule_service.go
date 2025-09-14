package service

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
)

// AccessControlRuleService defines the interface for access control rule business logic.
type AccessControlRuleService interface {
	GetAll(searchQuery schema.AccessControlRuleSearchQuery) ([]model.AccessControlRule, error)
	GetByID(id string) (*model.AccessControlRule, error)
	Save(id string, ruleModel *model.AccessControlRule) error
	Delete(id string) error
}

type accessControlRuleServiceImpl struct {
	accessControlRuleRepo repository.AccessControlRuleRepository
}

// NewAccessControlRuleService creates a new instance of AccessControlRuleService.
func NewAccessControlRuleService(accessControlRuleRepo repository.AccessControlRuleRepository) AccessControlRuleService {
	return &accessControlRuleServiceImpl{
		accessControlRuleRepo: accessControlRuleRepo,
	}
}

// GetAll retrieves all access control rules.
func (s *accessControlRuleServiceImpl) GetAll(searchQuery schema.AccessControlRuleSearchQuery) ([]model.AccessControlRule, error) {
	return s.accessControlRuleRepo.GetAll(searchQuery)
}

// GetByID retrieves an access control rule by its ID.
func (s *accessControlRuleServiceImpl) GetByID(id string) (*model.AccessControlRule, error) {
	return s.accessControlRuleRepo.GetByID(id)
}

// Save creates or updates an access control rule.
func (s *accessControlRuleServiceImpl) Save(id string, ruleModel *model.AccessControlRule) error {
	// Validate
	if ruleModel.Name == "" {
		return fmt.Errorf("rule name cannot be empty")
	}

	if id == "" {
		// Check if name already exists for new rule creation
		isExist, err := s.accessControlRuleRepo.IsExistName(ruleModel.Name)
		if err != nil {
			return fmt.Errorf("failed to check rule name existence: %w", err)
		}
		if isExist {
			return fmt.Errorf("rule with name '%s' already exists", ruleModel.Name)
		}
		// Create new rule
		if err := s.accessControlRuleRepo.Create(ruleModel); err != nil {
			return fmt.Errorf("failed to create rule: %w", err)
		}
	} else {
		// Update existing rule
		if err := s.accessControlRuleRepo.Update(ruleModel); err != nil {
			return fmt.Errorf("failed to update rule: %w", err)
		}
	}

	return nil
}

// Delete deletes an access control rule by its ID.
func (s *accessControlRuleServiceImpl) Delete(id string) error {
	return s.accessControlRuleRepo.Delete(id)
}
