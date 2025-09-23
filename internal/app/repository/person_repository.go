package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// --- Interfaces ---

// PersonRepository is the interface for person data access.
type PersonRepository interface {
	GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error)
	GetByID(id uuid.UUID) (*model.Person, error)
	Create(person *model.Person) error
	Update(id string, person *model.Person) error
	Delete(id uuid.UUID) error
	IsExistPersonID(personID string, excludeID uuid.UUID) (bool, error)
	IsExistName(firstName string, lastName string, excludeID uuid.UUID) (bool, error)
}

// PersonCardRepository is the interface for person card data access.
type PersonCardRepository interface {
	Create(cards []model.PersonCard) error
	GetCardNumbersByPersonID(personID string) ([]string, error)
	DeleteByPersonID(personID string) error
}

// PersonLicensePlateRepository is the interface for person license plate data access.
type PersonLicensePlateRepository interface {
	Create(plates []model.PersonLicensePlate) error
	GetLicensePlateTextsByPersonID(personID string) ([]string, error)
	DeleteByPersonID(personID string) error
}

// --- Implementations ---

// personRepositoryImpl is the implementation of PersonRepository.
type personRepositoryImpl struct {
	db *gorm.DB
}

// personCardRepositoryImpl is the implementation of PersonCardRepository.
type personCardRepositoryImpl struct {
	db *gorm.DB
}

// personLicensePlateRepositoryImpl is the implementation of PersonLicensePlateRepository.
type personLicensePlateRepositoryImpl struct {
	db *gorm.DB
}

// --- Constructors ---

// NewPersonRepository creates a new instance of PersonRepository.
func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepositoryImpl{db: db}
}

// NewPersonCardRepository creates a new instance of PersonCardRepository.
func NewPersonCardRepository(db *gorm.DB) PersonCardRepository {
	return &personCardRepositoryImpl{db: db}
}

// NewPersonLicensePlateRepository creates a new instance of PersonLicensePlateRepository.
func NewPersonLicensePlateRepository(db *gorm.DB) PersonLicensePlateRepository {
	return &personLicensePlateRepositoryImpl{db: db}
}

// --- PersonRepository Methods ---

// GetAll retrieves all persons with pagination.
func (r *personRepositoryImpl) GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error) {
	var persons []model.Person
	query := r.db.Model(&model.Person{})

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

	if !searchQuery.All {
		offset := (searchQuery.Page - 1) * searchQuery.Limit
		query = query.Offset(offset).Limit(searchQuery.Limit)
	}

	if err := query.Find(&persons).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve persons: %w", err)
	}

	return persons, nil
}

// GetByID retrieves a person by its ID.
func (r *personRepositoryImpl) GetByID(id uuid.UUID) (*model.Person, error) {
	var person model.Person
	if err := r.db.First(&person, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

// Create creates a new person record.
func (r *personRepositoryImpl) Create(person *model.Person) error {
	return r.db.Create(person).Error
}

// Update updates an existing person record.
func (r *personRepositoryImpl) Update(id string, person *model.Person) error {
	return r.db.Model(&model.Person{}).Where("id = ?", id).Updates(person).Error
}

// Delete deletes a person by its ID and all related records.
func (r *personRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create repository instances with the transaction
		txCardRepo := NewPersonCardRepository(tx)
		txLicenseRepo := NewPersonLicensePlateRepository(tx)

		// Delete related records first
		if err := txCardRepo.DeleteByPersonID(id.String()); err != nil {
			return err
		}
		if err := txLicenseRepo.DeleteByPersonID(id.String()); err != nil {
			return err
		}

		// Delete the person record itself
		if err := tx.Unscoped().Where("id = ?", id).Delete(&model.Person{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// IsExistPersonID checks if a person with the given PersonID exists.
func (r *personRepositoryImpl) IsExistPersonID(personID string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.Person{}).Where("person_id = ? AND deleted_at IS NULL", strings.ToLower(personID))
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check person ID existence: %w", err)
	}
	return count > 0, nil
}

// IsExistName checks if a person with the given first name and last name exists.
func (r *personRepositoryImpl) IsExistName(firstName string, lastName string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.Person{}).Where("first_name = ? AND last_name = ? AND deleted_at IS NULL", firstName, lastName)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check person name existence: %w", err)
	}
	return count > 0, nil
}

// --- PersonCardRepository Methods ---

// Create inserts multiple PersonCard records.
func (r *personCardRepositoryImpl) Create(cards []model.PersonCard) error {
	if len(cards) == 0 {
		return nil
	}
	return r.db.Create(&cards).Error
}

// GetCardNumbersByPersonID retrieves card numbers for a person.
func (r *personCardRepositoryImpl) GetCardNumbersByPersonID(personID string) ([]string, error) {
	var cardNumbers []string
	err := r.db.Model(&model.PersonCard{}).
		Select("card_number").
		Where("person_id = ?", personID).
		Find(&cardNumbers).Error
	return cardNumbers, err
}

// DeleteByPersonID deletes all person cards for a given person ID.
func (r *personCardRepositoryImpl) DeleteByPersonID(personID string) error {
	return r.db.Unscoped().Where("person_id = ?", personID).Delete(&model.PersonCard{}).Error
}

// --- PersonLicensePlateRepository Methods ---

// Create inserts multiple PersonLicensePlate records.
func (r *personLicensePlateRepositoryImpl) Create(plates []model.PersonLicensePlate) error {
	if len(plates) == 0 {
		return nil
	}
	return r.db.Create(&plates).Error
}

// GetLicensePlateTextsByPersonID retrieves license plate texts for a person.
func (r *personLicensePlateRepositoryImpl) GetLicensePlateTextsByPersonID(personID string) ([]string, error) {
	var licensePlates []string
	err := r.db.Model(&model.PersonLicensePlate{}).
		Select("license_plate_text").
		Where("person_id = ?", personID).
		Find(&licensePlates).Error
	return licensePlates, err
}

// DeleteByPersonID deletes all person license plates for a given person ID.
func (r *personLicensePlateRepositoryImpl) DeleteByPersonID(personID string) error {
	return r.db.Unscoped().Where("person_id = ?", personID).Delete(&model.PersonLicensePlate{}).Error
}
