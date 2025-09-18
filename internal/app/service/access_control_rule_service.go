package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlRuleService defines the interface for access control rule business logic.
type AccessControlRuleService interface {
	GetAll(searchQuery schema.AccessControlRuleSearchQuery) ([]model.AccessControlRule, error)
	GetByID(id string) (*model.AccessControlRule, error)
	Create(bodyRequest *schema.AccessControlRuleRequest) (*model.AccessControlRule, error)
	Update(id string, bodyRequest *schema.AccessControlRuleRequest) (*model.AccessControlRule, error)
	PartialUpdate(id string, bodyRequest *schema.AccessControlRuleRequest) (*model.AccessControlRule, error)
	Delete(id string) error
	ConvertToResponse(ruleModel *model.AccessControlRule) (*schema.AccessControlRuleResponse, error)
	GetGroupsInfo(groupIDs []string) ([]schema.AccessControlGroupInfoResponse, error)
}

type accessControlRuleServiceImpl struct {
	accessControlRuleRepo  repository.AccessControlRuleRepository
	accessControlGroupRepo repository.AccessControlGroupRepository // ต้องใช้สำหรับตรวจสอบ Group ID และดึงข้อมูล
	db                     *gorm.DB
}

// NewAccessControlRuleService creates a new instance of AccessControlRuleService.
func NewAccessControlRuleService(accessControlRuleRepo repository.AccessControlRuleRepository, accessControlGroupRepo repository.AccessControlGroupRepository, db *gorm.DB) AccessControlRuleService {
	return &accessControlRuleServiceImpl{
		accessControlRuleRepo:  accessControlRuleRepo,
		accessControlGroupRepo: accessControlGroupRepo,
		db:                     db,
	}
}

// GetAll retrieves all access control rules.
func (s *accessControlRuleServiceImpl) GetAll(searchQuery schema.AccessControlRuleSearchQuery) ([]model.AccessControlRule, error) {
	return s.accessControlRuleRepo.GetAll(searchQuery)
}

// GetByID retrieves an access control rule by its ID.
func (s *accessControlRuleServiceImpl) GetByID(id string) (*model.AccessControlRule, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.accessControlRuleRepo.GetByID(idUUID)
}

// Create creates a new access control rule.
func (s *accessControlRuleServiceImpl) Create(bodyRequest *schema.AccessControlRuleRequest) (*model.AccessControlRule, error) {

	// Set default value and Validate
	bodyRequest, err := s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	// Validate (check duplicates)
	if err := s.validateBodyRequest(*bodyRequest, nil); err != nil {
		return nil, err
	}

	ruleModel := &model.AccessControlRule{
		Name: bodyRequest.Name,
	}

	// ใช้ Transaction เพื่อให้แน่ใจว่าทั้ง Rule และ Group ถูกสร้างหรือยกเลิกพร้อมกัน
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAccessControlRuleRepository(tx) // ต้องสร้าง NewAccessControlRuleRepository ที่รับ tx
		// 1. Create Rule
		if err := txRepo.Create(ruleModel); err != nil {
			return fmt.Errorf("failed to create rule: %w", err)
		}
		// 2. Create Rule Groups
		if len(bodyRequest.AccessControlGroupIDs) > 0 {
			ruleGroups, err := s.createRuleGroupModels(ruleModel.ID.String(), bodyRequest.AccessControlGroupIDs)
			if err != nil {
				return err
			}
			// Add rule group
			if err := txRepo.CreateRuleGroups(ruleGroups); err != nil {
				return fmt.Errorf("failed to create rule groups: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return ruleModel, nil
}

// Update updates an existing access control rule (Full Replacement).
func (s *accessControlRuleServiceImpl) Update(id string, bodyRequest *schema.AccessControlRuleRequest) (*model.AccessControlRule, error) {

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	ruleModel, err := s.accessControlRuleRepo.GetByID(idUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing rule: %w", err)
	}
	if ruleModel == nil {
		return nil, fmt.Errorf("rule with ID '%s' not found", id)
	}

	// Set default value (เพื่อให้แน่ใจว่า GroupIDs ถูกเคลียร์ถ้าไม่ได้ส่งมา)
	bodyRequest, err = s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	// Validate (check duplicates)
	if err := s.validateBodyRequest(*bodyRequest, ruleModel); err != nil {
		return nil, err
	}

	// Update model
	ruleModel.Name = bodyRequest.Name

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAccessControlRuleRepository(tx)
		// 1. Update Rule
		if err := txRepo.Update(ruleModel); err != nil {
			return fmt.Errorf("failed to update rule: %w", err)
		}

		// 2. Delete Rule Groups all & Recreate
		if err := txRepo.DeleteRuleGroupsByRuleID(idUUID, tx); err != nil {
			return fmt.Errorf("failed to delete old rule groups: %w", err)
		}

		// Create Rule Groups
		if len(bodyRequest.AccessControlGroupIDs) > 0 {
			ruleGroups, err := s.createRuleGroupModels(id, bodyRequest.AccessControlGroupIDs)
			if err != nil {
				return err
			}
			if err := txRepo.CreateRuleGroups(ruleGroups); err != nil {
				return fmt.Errorf("failed to create new rule groups: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ruleModel, nil
}

// PartialUpdate performs a partial update on an existing access control rule.
func (s *accessControlRuleServiceImpl) PartialUpdate(id string, bodyRequest *schema.AccessControlRuleRequest) (*model.AccessControlRule, error) {

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	ruleModel, err := s.accessControlRuleRepo.GetByID(idUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing rule: %w", err)
	}
	if ruleModel == nil {
		return nil, fmt.Errorf("rule with ID '%s' not found", id)
	}

	// Validate (check duplicates)
	// Note: ไม่เรียก validateAndSetDefaultValues เพื่อให้ GroupIDs เป็น nil ถ้าไม่ได้ส่งมา
	if err := s.validateBodyRequest(*bodyRequest, ruleModel); err != nil {
		return nil, err
	}

	// Update model
	// Note: Name ใน Request เป็น string ธรรมดา ดังนั้นต้องเช็กว่า Name ถูกส่งมาหรือไม่
	// แต่ใน Request Schema Name ถูกตั้งเป็น validate:"required" ดังนั้นจึงต้องมีค่าเสมอ
	if bodyRequest.Name != "" {
		ruleModel.Name = bodyRequest.Name
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAccessControlRuleRepository(tx)

		// 1. Update Rule
		if err := txRepo.Update(ruleModel); err != nil {
			return fmt.Errorf("failed to update rule: %w", err)
		}

		// 2. Update Rule Groups (เฉพาะถ้ามีการส่ง AccessControlGroupIDs มาใน body request)
		if bodyRequest.AccessControlGroupIDs != nil {
			// Delete Rule Groups all
			if err := txRepo.DeleteRuleGroupsByRuleID(idUUID, tx); err != nil {
				return fmt.Errorf("failed to delete old rule groups: %w", err)
			}
			// Create Rule Groups
			if len(bodyRequest.AccessControlGroupIDs) > 0 {
				ruleGroups, err := s.createRuleGroupModels(id, bodyRequest.AccessControlGroupIDs)
				if err != nil {
					return err
				}
				if err := txRepo.CreateRuleGroups(ruleGroups); err != nil {
					return fmt.Errorf("failed to create new rule groups: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ruleModel, nil
}

// Delete deletes an access control rule by its ID.
func (s *accessControlRuleServiceImpl) Delete(id string) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}
	_, err = s.accessControlRuleRepo.GetByID(idUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("rule with ID '%s' not found", id)
		}
	}
	// Note: การลบ Transaction ถูกจัดการภายใน AccessControlRuleRepository.Delete แล้ว
	return s.accessControlRuleRepo.Delete(idUUID)
}

// ConvertToResponse converts a rule model to a response schema.
func (s *accessControlRuleServiceImpl) ConvertToResponse(ruleModel *model.AccessControlRule) (*schema.AccessControlRuleResponse, error) {

	// 1. ดึง Group IDs ที่ผูกกับ Rule นี้
	groupIDs, err := s.accessControlRuleRepo.GetGroupIDsByRuleID(ruleModel.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group IDs for rule: %w", err)
	}

	// 2. ดึงข้อมูล Group เต็ม
	groupResponses, err := s.GetGroupsInfo(groupIDs)
	if err != nil {
		return nil, err
	}

	// 3. สร้าง Response
	response := &schema.AccessControlRuleResponse{
		ID:                  ruleModel.ID.String(),
		Name:                ruleModel.Name,
		AccessControlGroups: groupResponses,
	}

	return response, nil
}

// ----------> RELATE FUNCTION <-----------------------//

// GetGroupsInfo retrieves full group info from group IDs.
func (s *accessControlRuleServiceImpl) GetGroupsInfo(groupIDs []string) ([]schema.AccessControlGroupInfoResponse, error) {
	if len(groupIDs) == 0 {
		return []schema.AccessControlGroupInfoResponse{}, nil
	}

	var groupsInfo []schema.AccessControlGroupInfoResponse

	for _, idStr := range groupIDs {
		idUUID, err := uuid.Parse(idStr)
		if err != nil {
			// อาจจะ log warning แทนการ return error หาก ID ไม่ถูกต้อง
			continue
		}
		// ใช้ accessControlGroupRepo.GetByID เพื่อดึงข้อมูล Group เต็ม
		group, err := s.accessControlGroupRepo.GetByID(idUUID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Group ไม่พบแล้ว ข้ามไป
				continue
			}
			return nil, fmt.Errorf("failed to get group by ID '%s': %w", idStr, err)
		}
		groupsInfo = append(groupsInfo, schema.AccessControlGroupInfoResponse{
			ID:   group.ID.String(),
			Name: group.Name,
		})
	}
	return groupsInfo, nil
}

// ----------> INNER FUNCTION <-----------------------//

// createRuleGroupModels converts group IDs to AccessControlRuleGroup models.
func (s *accessControlRuleServiceImpl) createRuleGroupModels(ruleID string, groupIDs []string) ([]model.AccessControlRuleGroup, error) {
	var ruleGroups []model.AccessControlRuleGroup
	for _, groupID := range groupIDs {
		// ตรวจสอบว่า Group ID นั้นมีอยู่จริงหรือไม่ (แนะนำ)
		groupUUID, err := uuid.Parse(groupID)
		if err != nil {
			return nil, fmt.Errorf("invalid group ID format: %s", groupID)
		}
		_, err = s.accessControlGroupRepo.GetByID(groupUUID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("group with ID '%s' not found", groupID)
			}
			return nil, fmt.Errorf("failed to check group existence: %w", err)
		}

		ruleGroups = append(ruleGroups, model.AccessControlRuleGroup{
			AccessControlRuleID:  ruleID,
			AccessControlGroupID: groupID,
		})
	}
	return ruleGroups, nil
}

// Validate and set default
func (s *accessControlRuleServiceImpl) validateAndSetDefaultValues(bodyRequest *schema.AccessControlRuleRequest) (*schema.AccessControlRuleRequest, error) {
	// Name เป็น string ธรรมดาและ required ใน Request Schema จึงไม่ต้องเช็ก nil
	if bodyRequest.Name == "" {
		// แม้ว่า validate:"required" จะช่วย แต่การเช็กใน service ก็ดีกว่า
		return nil, fmt.Errorf("rule name cannot be empty")
	}
	// ถ้า AccessControlGroupIDs เป็น nil ให้ตั้งเป็น slice ว่าง
	if bodyRequest.AccessControlGroupIDs == nil {
		bodyRequest.AccessControlGroupIDs = []string{}
	}
	return bodyRequest, nil
}

// validateBodyRequest checks for duplicate rule name.
func (s *accessControlRuleServiceImpl) validateBodyRequest(bodyRequest schema.AccessControlRuleRequest, ruleModel *model.AccessControlRule) error {

	excludeID := uuid.Nil
	if ruleModel != nil {
		excludeID = ruleModel.ID
	}

	// Check duplicate name
	// Name เป็น string ธรรมดาและ required
	isExistName, err := s.accessControlRuleRepo.IsExistName(bodyRequest.Name, excludeID)
	if err != nil {
		return fmt.Errorf("failed to check rule name existence: %w", err)
	}
	if isExistName {
		return fmt.Errorf("accessControlRule name is already exist: %s", bodyRequest.Name)
	}

	// Note: การตรวจสอบว่า AccessControlGroupIDs มีอยู่จริงหรือไม่ ถูกย้ายไปทำใน createRuleGroupModels
	return nil
}
