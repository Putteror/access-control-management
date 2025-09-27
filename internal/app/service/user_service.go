package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// UserService defines the interface for user business logic.
type UserService interface {
	GetAll(searchQuery schema.UserSearchQuery) ([]model.User, error)
	GetByID(id string) (*model.User, error)
	Create(bodyRequest *schema.UserRequest) (*model.User, error)
	Update(id string, bodyRequest *schema.UserRequest) (*model.User, error)
	PartialUpdate(id string, bodyRequest *schema.UserRequest) (*model.User, error)
	Delete(id string) error
	ConvertToResponse(userModel *model.User) (*schema.UserResponse, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
	db       *gorm.DB
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo repository.UserRepository, db *gorm.DB) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
		db:       db,
	}
}

// -----------------------------------------------------------------------------
// --- CRUD Operations ---
// -----------------------------------------------------------------------------

// GetAll retrieves all user records.
func (s *userServiceImpl) GetAll(searchQuery schema.UserSearchQuery) ([]model.User, error) {
	return s.userRepo.GetAll(searchQuery)
}

// GetByID retrieves a user record by its ID.
func (s *userServiceImpl) GetByID(id string) (*model.User, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	userModel, err := s.userRepo.GetByID(idUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID '%s' not found", id)
		}
		return nil, err
	}
	return userModel, nil
}

// Create creates a new user and their permission.
func (s *userServiceImpl) Create(bodyRequest *schema.UserRequest) (*model.User, error) {

	// 1. Validate mandatory fields (Username, Password, Permission)
	if err := s.validateMandatoryFields(bodyRequest); err != nil {
		return nil, err
	}

	// 2. Validate business rules (Check duplicate username)
	if err := s.validateBodyRequest(bodyRequest, uuid.Nil); err != nil {
		return nil, err
	}

	// 3. Hash Password (Mock)
	passwordHash, err := common.HashPassword(*bodyRequest.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. Create Models
	permissionModel := s.createPermissionModel(bodyRequest.Permission)
	userModel := &model.User{
		Username:     *bodyRequest.Username,
		PasswordHash: passwordHash,
		Status:       s.getOrDefaultStatus(bodyRequest.Status),
	}

	// 5. Create User and Permission in a single transaction
	if err := s.userRepo.CreateUserWithPermission(userModel, permissionModel); err != nil {
		return nil, fmt.Errorf("failed to create user and permission: %w", err)
	}

	// 6. Set relationship for response mapping
	userModel.Permission = *permissionModel

	return userModel, nil
}

// Update updates an existing user and their permission (Full Replacement).
func (s *userServiceImpl) Update(id string, bodyRequest *schema.UserRequest) (*model.User, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	// 1. Get existing user model (Preloaded with Permission)
	userModel, err := s.userRepo.GetByID(idUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to get existing user: %w", err)
	}

	// 2. Validate mandatory fields (Username, Permission)
	if err := s.validateMandatoryFields(bodyRequest); err != nil {
		return nil, err
	}

	// 3. Validate business rules (Check duplicate username, excluding itself)
	if err := s.validateBodyRequest(bodyRequest, idUUID); err != nil {
		return nil, err
	}

	// 4. Update User Model
	userModel.Username = *bodyRequest.Username
	userModel.Status = s.getOrDefaultStatus(bodyRequest.Status)

	// Update Password only if a new one is provided
	if bodyRequest.Password != nil && *bodyRequest.Password != "" {
		passwordHash, err := common.HashPassword(*bodyRequest.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		userModel.PasswordHash = passwordHash
	}

	// 5. Update Permission Model (ใช้ ID เดิม)
	permissionModel := &userModel.Permission
	s.mapPermissionRequestToModel(bodyRequest.Permission, permissionModel)

	// 6. Update User and Permission in a single transaction
	if err := s.userRepo.UpdateUserAndPermission(userModel, permissionModel); err != nil {
		return nil, fmt.Errorf("failed to update user and permission: %w", err)
	}

	return userModel, nil
}

// PartialUpdate performs a partial update on an existing user record.
func (s *userServiceImpl) PartialUpdate(id string, bodyRequest *schema.UserRequest) (*model.User, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	// 1. Get existing user model (Preloaded with Permission)
	userModel, err := s.userRepo.GetByID(idUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to get existing user: %w", err)
	}

	// 2. Validate business rules (Check duplicate username, excluding itself)
	if err := s.validateBodyRequest(bodyRequest, idUUID); err != nil {
		return nil, err
	}

	// 3. Apply Partial Updates to User Model
	if bodyRequest.Username != nil && *bodyRequest.Username != "" {
		userModel.Username = *bodyRequest.Username
	}
	if bodyRequest.Status != nil && *bodyRequest.Status != "" {
		userModel.Status = *bodyRequest.Status
	}
	if bodyRequest.Password != nil && *bodyRequest.Password != "" {
		passwordHash, err := common.HashPassword(*bodyRequest.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		userModel.PasswordHash = passwordHash
	}

	// 4. Apply Partial Updates to Permission Model (ถ้ามีการส่ง Permission มา)
	if bodyRequest.Permission != nil {
		// ใช้ ID เดิม
		permissionModel := &userModel.Permission
		s.mapPermissionRequestToModel(bodyRequest.Permission, permissionModel)

		// 5. Update User and Permission in a single transaction
		if err := s.userRepo.UpdateUserAndPermission(userModel, permissionModel); err != nil {
			return nil, fmt.Errorf("failed to update user and permission: %w", err)
		}
	} else {
		// 5.1 ถ้าไม่มีการส่ง Permission มา ให้อัพเดทเฉพาะ User
		if err := s.userRepo.Update(userModel); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return userModel, nil
}

// Delete deletes a user record by its ID.
func (s *userServiceImpl) Delete(id string) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}

	// ตรวจสอบการมีอยู่ก่อนลบ
	_, err = s.userRepo.GetByID(idUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get user by ID: %w", err)
	}

	// การลบ Transaction ถูกจัดการภายใน UserRepository.Delete แล้ว
	return s.userRepo.Delete(idUUID)
}

// ConvertToResponse converts a user model to a response schema.
func (s *userServiceImpl) ConvertToResponse(userModel *model.User) (*schema.UserResponse, error) {

	// สร้าง UserPermissionResponse จาก UserPermission model
	permissionResponse := schema.UserPermissionResponse{
		ID:                       userModel.Permission.ID.String(),
		PeoplePermission:         userModel.Permission.PeoplePermission,
		DevicePermission:         userModel.Permission.DevicePermission,
		RulePermission:           userModel.Permission.RulePermission,
		TimeAttendancePermission: userModel.Permission.TimeAttendancePermission,
		ReportPermission:         userModel.Permission.ReportPermission,
		NotificationPermission:   userModel.Permission.NotificationPermission,
		SystemLogPermission:      userModel.Permission.SystemLogPermission,
	}

	// สร้าง UserResponse
	response := &schema.UserResponse{
		ID:         userModel.ID.String(),
		Username:   userModel.Username,
		Status:     userModel.Status,
		Permission: permissionResponse,
	}

	return response, nil
}

// -----------------------------------------------------------------------------
// --- Inner Functions ---
// -----------------------------------------------------------------------------

// getOrDefaultStatus returns the status from request or a default value (e.g., "active").
func (s *userServiceImpl) getOrDefaultStatus(statusPtr *string) string {
	if statusPtr != nil && *statusPtr != "" {
		// Note: ไม่มี common.DefaultUserStatus, ใช้ "active" เป็นค่า default ชั่วคราว
		return *statusPtr
	}
	return "active"
}

// createPermissionModel converts UserPermissionRequest to UserPermission model, setting defaults to false.
func (s *userServiceImpl) createPermissionModel(request *schema.UserPermissionRequest) *model.UserPermission {
	// กำหนดค่าเริ่มต้นเป็น false ทั้งหมด
	permissionModel := &model.UserPermission{
		PeoplePermission:         false,
		DevicePermission:         false,
		RulePermission:           false,
		TimeAttendancePermission: false,
		ReportPermission:         false,
		NotificationPermission:   false,
		SystemLogPermission:      false,
	}

	// Map ค่าจาก Request ที่เป็น Pointer (*bool)
	s.mapPermissionRequestToModel(request, permissionModel)

	return permissionModel
}

// mapPermissionRequestToModel applies non-nil fields from request to model.
func (s *userServiceImpl) mapPermissionRequestToModel(request *schema.UserPermissionRequest, model *model.UserPermission) {
	if request.PeoplePermission != nil {
		model.PeoplePermission = *request.PeoplePermission
	}
	if request.DevicePermission != nil {
		model.DevicePermission = *request.DevicePermission
	}
	if request.RulePermission != nil {
		model.RulePermission = *request.RulePermission
	}
	if request.TimeAttendancePermission != nil {
		model.TimeAttendancePermission = *request.TimeAttendancePermission
	}
	if request.ReportPermission != nil {
		model.ReportPermission = *request.ReportPermission
	}
	if request.NotificationPermission != nil {
		model.NotificationPermission = *request.NotificationPermission
	}
	if request.SystemLogPermission != nil {
		model.SystemLogPermission = *request.SystemLogPermission
	}
}

// validateMandatoryFields checks for required fields for Create/Update.
func (s *userServiceImpl) validateMandatoryFields(bodyRequest *schema.UserRequest) error {
	if bodyRequest.Username == nil || *bodyRequest.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// สำหรับ Create, Password ต้องมี
	if bodyRequest.Password == nil || *bodyRequest.Password == "" {
		// ในกรณี Update/PartialUpdate จะไม่ตรวจสอบ
	}

	if bodyRequest.Permission == nil {
		return fmt.Errorf("permission configuration is mandatory")
	}
	return nil
}

// validateBodyRequest checks for duplicate username.
func (s *userServiceImpl) validateBodyRequest(bodyRequest *schema.UserRequest, excludeID uuid.UUID) error {

	// Check duplicate name only if Username is provided (for PartialUpdate)
	if bodyRequest.Username != nil && *bodyRequest.Username != "" {
		isExistUsername, err := s.userRepo.IsExistUsername(*bodyRequest.Username, excludeID)
		if err != nil {
			return fmt.Errorf("failed to check username existence: %w", err)
		}

		if isExistUsername {
			// ใช้ข้อความที่บันทึกไว้ตามคำขอของผู้ใช้ (2025-07-05)
			return fmt.Errorf("ApplicationForm name is already exist")
		}
	}

	return nil
}
