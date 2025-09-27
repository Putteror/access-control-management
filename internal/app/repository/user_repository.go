package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// UserRepository is the interface for user data access.
type UserRepository interface {
	// User CRUD operations
	GetAll(searchQuery schema.UserSearchQuery) ([]model.User, error)
	GetByID(id uuid.UUID) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uuid.UUID) error
	IsExistUsername(username string, excludeID uuid.UUID) (bool, error)

	// UserPermission relationship operations (if needed separately)
	GetPermissionByID(id string) (*model.UserPermission, error) // ID ของ UserPermission
	CreatePermission(permission *model.UserPermission, tx *gorm.DB) error
	UpdatePermission(permission *model.UserPermission, tx *gorm.DB) error
	DeletePermission(id string, tx *gorm.DB) error

	// Transactional operations
	CreateUserWithPermission(user *model.User, permission *model.UserPermission) error
	UpdateUserAndPermission(user *model.User, permission *model.UserPermission) error
}

// userRepositoryImpl is the implementation of UserRepository.
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// -----------------------------------------------------------------------------
// --- User CRUD Operations ---
// -----------------------------------------------------------------------------

// GetAll retrieves all user records with pagination and search.
func (r *userRepositoryImpl) GetAll(searchQuery schema.UserSearchQuery) ([]model.User, error) {
	var users []model.User

	// Preload "Permission" relationship
	query := r.db.Model(&model.User{}).Preload("Permission")

	// Search filters
	if searchQuery.Username != "" {
		query = query.Where("username ILIKE ?", "%"+searchQuery.Username+"%")
	}
	if searchQuery.Status != "" {
		query = query.Where("status = ?", searchQuery.Status)
	}

	// Pagination
	page := searchQuery.Page
	limit := searchQuery.Limit
	offset := (page - 1) * limit

	if page > 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated user records: %w", err)
	}

	return users, nil
}

// GetByID retrieves a user record by its ID, including permission details.
func (r *userRepositoryImpl) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	// Preload "Permission" relationship
	if err := r.db.Preload("Permission").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user record by their username.
func (r *userRepositoryImpl) GetByUsername(username string) (*model.User, error) {
	var user model.User
	// ไม่จำเป็นต้อง Preload Permission ในกรณีนี้ถ้าใช้เพื่อการ Login/Auth เพียงอย่างเดียว
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user record.
// Note: ควรใช้ CreateUserWithPermission สำหรับการสร้าง User พร้อมสิทธิ์
func (r *userRepositoryImpl) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update updates an existing user record.
func (r *userRepositoryImpl) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user record by its ID.
// Note: เราจะใช้ Transaction เพื่อลบ User และ Permission ที่เกี่ยวข้อง
func (r *userRepositoryImpl) Delete(id uuid.UUID) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. ดึงข้อมูล User เพื่อหา PermissionID
		var user model.User
		if err := tx.First(&user, "id = ?", id).Error; err != nil {
			return err
		}

		// 2. ลบ UserPermission
		if err := r.DeletePermission(user.PermissionID, tx); err != nil {
			return err
		}

		// 3. ลบ User
		if err := tx.Unscoped().Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// IsExistUsername checks if a user record with the given username exists.
func (r *userRepositoryImpl) IsExistUsername(username string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.User{}).Where("username = ? AND deleted_at IS NULL", username)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

// -----------------------------------------------------------------------------
// --- UserPermission Relationship Operations ---
// -----------------------------------------------------------------------------

// GetPermissionByID retrieves a UserPermission record by its ID.
func (r *userRepositoryImpl) GetPermissionByID(id string) (*model.UserPermission, error) {
	var permission model.UserPermission
	if err := r.db.First(&permission, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// CreatePermission creates a new UserPermission record within the provided transaction.
func (r *userRepositoryImpl) CreatePermission(permission *model.UserPermission, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(permission).Error
}

// UpdatePermission updates an existing UserPermission record within the provided transaction.
func (r *userRepositoryImpl) UpdatePermission(permission *model.UserPermission, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(permission).Error
}

// DeletePermission deletes a UserPermission record by its ID within the provided transaction.
func (r *userRepositoryImpl) DeletePermission(id string, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Unscoped().Where("id = ?", id).Delete(&model.UserPermission{}).Error
}

// -----------------------------------------------------------------------------
// --- Transactional Operations ---
// -----------------------------------------------------------------------------

// CreateUserWithPermission creates both a User and their associated UserPermission in a single transaction.
func (r *userRepositoryImpl) CreateUserWithPermission(user *model.User, permission *model.UserPermission) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. สร้าง Permission
		if err := r.CreatePermission(permission, tx); err != nil {
			return err
		}

		user.PermissionID = permission.ID.String()

		return tx.Create(user).Error
	})
}

// UpdateUserAndPermission updates both User and UserPermission in a single transaction.
func (r *userRepositoryImpl) UpdateUserAndPermission(user *model.User, permission *model.UserPermission) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. อัพเดท Permission
		if err := r.UpdatePermission(permission, tx); err != nil {
			return err
		}

		// 2. อัพเดท User
		return tx.Save(user).Error
	})
}
