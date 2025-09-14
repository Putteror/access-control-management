package repository

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"

	"gorm.io/gorm"
)

type PersonRepository interface {
	GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error)
	GetByID(id string) (*model.Person, error)
	Create(people *model.Person) error
	Update(people *model.Person) error
	Delete(id string) error
}

type personRepositoryImpl struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepositoryImpl{db: db}
}

func (r *personRepositoryImpl) FindAll() ([]model.Person, error) {
	var people []model.Person
	if err := r.db.Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

func (r *personRepositoryImpl) GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error) {
	var persons []model.Person

	query := r.db.Model(&model.Person{})

	// Add search conditions based on the provided query
	if searchQuery.FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+searchQuery.FirstName+"%")
	}
	if searchQuery.LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+searchQuery.LastName+"%")
	}
	if searchQuery.Company != "" {
		query = query.Where("company ILIKE ?", "%"+searchQuery.Company+"%")
	}
	if searchQuery.Department != "" {
		query = query.Where("department ILIKE ?", "%"+searchQuery.Department+"%")
	}
	if searchQuery.JobPosition != "" {
		query = query.Where("job_position ILIKE ?", "%"+searchQuery.JobPosition+"%")
	}
	if searchQuery.MobileNumber != "" {
		query = query.Where("mobile_number ILIKE ?", "%"+searchQuery.MobileNumber+"%")
	}
	if searchQuery.Email != "" {
		query = query.Where("email ILIKE ?", "%"+searchQuery.Email+"%")
	}

	var page int = searchQuery.Page
	var limit int = searchQuery.Limit
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&persons).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated persons: %w", err)
	}
	return persons, nil
}

func (r *personRepositoryImpl) GetByID(id string) (*model.Person, error) {
	var person model.Person
	if err := r.db.First(&person, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *personRepositoryImpl) Create(person *model.Person) error {
	return r.db.Create(person).Error
}

func (r *personRepositoryImpl) Update(person *model.Person) error {
	return r.db.Save(person).Error
}

func (r *personRepositoryImpl) Delete(id string) error {
	return r.db.Unscoped().Where("id = ?", id).Delete(&model.Person{}).Error
}

// IsExistName
func (r *personRepositoryImpl) IsExistName(name string) (bool, error) {
	var count int64
	var person model.Person
	result := r.db.Model(&person).Where("first_name = ? AND last_name = ? AND deleted_at IS NULL", name).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check rule existence: %w", result.Error)
	}
	return count > 0, nil
}
