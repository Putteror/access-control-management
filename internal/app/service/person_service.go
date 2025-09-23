package service

import (
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// PersonService defines the interface for person business logic.
type PersonService interface {
	GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error)
	GetByID(id string) (*model.Person, error)
	Save(id string, person *model.Person, faceImageFile *multipart.FileHeader, cardIDs []string, licensePlateTexts []string) error
	PartialUpdate(id string, person *model.Person, faceImageFile *multipart.FileHeader, cardIDs []string, licensePlateTexts []string) error
	Delete(id string) error
	ConvertToResponse(personModel *model.Person) (*schema.PersonResponse, error)
}

type personServiceImpl struct {
	personRepo         repository.PersonRepository
	personCardRepo     repository.PersonCardRepository
	personLicenseRepo  repository.PersonLicensePlateRepository
	accessRuleRepo     repository.AccessControlRuleRepository
	timeAttendanceRepo repository.AttendanceRepository
	db                 *gorm.DB
}

// NewPersonService creates a new instance of PersonService.
func NewPersonService(personRepo repository.PersonRepository, personCardRepo repository.PersonCardRepository, personLicenseRepo repository.PersonLicensePlateRepository, accessRuleRepo repository.AccessControlRuleRepository, timeAttendanceRepo repository.AttendanceRepository, db *gorm.DB) PersonService {
	return &personServiceImpl{
		personRepo:         personRepo,
		personCardRepo:     personCardRepo,
		personLicenseRepo:  personLicenseRepo,
		accessRuleRepo:     accessRuleRepo,
		timeAttendanceRepo: timeAttendanceRepo,
		db:                 db,
	}
}

// GetAll retrieves all persons.
func (s *personServiceImpl) GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error) {
	return s.personRepo.GetAll(searchQuery)
}

// GetByID retrieves a person by its ID.
func (s *personServiceImpl) GetByID(id string) (*model.Person, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.personRepo.GetByID(idUUID)
}

// Save creates or updates a person.
func (s *personServiceImpl) Save(id string, person *model.Person, faceImageFile *multipart.FileHeader, cardIDs []string, licensePlateTexts []string) error {
	isCreate := id == ""

	// Validate if PersonID and PersonName exist
	if err := s.validatePerson(isCreate, person); err != nil {
		return err
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txPersonRepo := repository.NewPersonRepository(tx)
		txCardRepo := repository.NewPersonCardRepository(tx)
		txLicenseRepo := repository.NewPersonLicensePlateRepository(tx)

		if isCreate {
			if err := txPersonRepo.Create(person); err != nil {
				return fmt.Errorf("failed to create person: %w", err)
			}
		} else {
			existingPerson, err := txPersonRepo.GetByID(uuid.MustParse(id))
			if err != nil {
				return fmt.Errorf("person with ID '%s' not found: %w", id, err)
			}
			if err := txPersonRepo.Update(id, person); err != nil {
				return fmt.Errorf("failed to update person: %w", err)
			}
			if faceImageFile != nil && existingPerson.FaceImagePath != nil && *existingPerson.FaceImagePath != "" {
				// TODO: Add logic to delete old image file from storage
				fmt.Printf("Deleting old image at path: %s\n", *existingPerson.FaceImagePath)
			}
		}

		// Handle relationships
		if err := s.handleRelationships(txCardRepo, txLicenseRepo, person.ID.String(), cardIDs, licensePlateTexts); err != nil {
			return err
		}

		if faceImageFile != nil {
			filePath, err := common.UploadFile(faceImageFile, person.ID.String())
			if err != nil {
				return fmt.Errorf("failed to upload face image: %w", err)
			}
			person.FaceImagePath = &filePath
			if err := txPersonRepo.Update(person.ID.String(), person); err != nil {
				return fmt.Errorf("failed to update person face image path: %w", err)
			}
		}

		return nil
	})
	return err
}

// PartialUpdate performs a partial update on a person.
func (s *personServiceImpl) PartialUpdate(id string, person *model.Person, faceImageFile *multipart.FileHeader, cardIDs []string, licensePlateTexts []string) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	existingPerson, err := s.personRepo.GetByID(idUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("person with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get existing person: %w", err)
	}

	// Validate
	if err := s.validatePerson(false, person); err != nil {
		return err
	}

	// Merge new data into existing model
	if person.FirstName != "" {
		existingPerson.FirstName = person.FirstName
	}
	if person.LastName != "" {
		existingPerson.LastName = person.LastName
	}
	if person.MiddleName != nil {
		existingPerson.MiddleName = person.MiddleName
	}
	if person.PersonType != "" {
		existingPerson.PersonType = person.PersonType
	}
	if person.PersonID != nil {
		existingPerson.PersonID = person.PersonID
	}
	if person.Gender != nil {
		existingPerson.Gender = person.Gender
	}
	if person.DateOfBirth != nil {
		existingPerson.DateOfBirth = person.DateOfBirth
	}
	if person.Company != nil {
		existingPerson.Company = person.Company
	}
	if person.Department != nil {
		existingPerson.Department = person.Department
	}
	if person.JobPosition != nil {
		existingPerson.JobPosition = person.JobPosition
	}
	if person.Address != nil {
		existingPerson.Address = person.Address
	}
	if person.MobileNumber != nil {
		existingPerson.MobileNumber = person.MobileNumber
	}
	if person.Email != nil {
		existingPerson.Email = person.Email
	}
	if person.IsVerified {
		existingPerson.IsVerified = person.IsVerified
	}
	if person.ActiveAt != nil {
		existingPerson.ActiveAt = person.ActiveAt
	}
	if person.ExpireAt != nil {
		existingPerson.ExpireAt = person.ExpireAt
	}
	if person.AccessControlRuleID != nil {
		existingPerson.AccessControlRuleID = person.AccessControlRuleID
	}
	if person.TimeAttendanceID != nil {
		existingPerson.TimeAttendanceID = person.TimeAttendanceID
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txPersonRepo := repository.NewPersonRepository(tx)
		txCardRepo := repository.NewPersonCardRepository(tx)
		txLicenseRepo := repository.NewPersonLicensePlateRepository(tx)

		// Handle face image upload/update
		if faceImageFile != nil {
			// TODO: Add logic to upload new image and set person.FaceImagePath
		}

		if err := txPersonRepo.Update(id, existingPerson); err != nil {
			return fmt.Errorf("failed to partial update person: %w", err)
		}

		// Handle relationships
		if cardIDs != nil || licensePlateTexts != nil {
			if err := s.handleRelationships(txCardRepo, txLicenseRepo, id, cardIDs, licensePlateTexts); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// Delete deletes a person by its ID.
func (s *personServiceImpl) Delete(id string) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}

	_, err = s.personRepo.GetByID(idUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("person with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get person by ID: %w", err)
	}

	return s.personRepo.Delete(idUUID)
}

// ConvertToResponse converts a person model to a response schema.
func (s *personServiceImpl) ConvertToResponse(personModel *model.Person) (*schema.PersonResponse, error) {
	var accessRule *schema.AccessControlRuleInfoResponse
	if personModel.AccessControlRuleID != nil && *personModel.AccessControlRuleID != "" {
		accessRuleModel, err := s.accessRuleRepo.GetByID(uuid.MustParse(*personModel.AccessControlRuleID))
		if err != nil {
			return nil, fmt.Errorf("failed to get access control rule: %w", err)
		}
		accessRule = &schema.AccessControlRuleInfoResponse{
			ID:   accessRuleModel.ID.String(),
			Name: accessRuleModel.Name,
		}
	}

	var timeAttendance *schema.AttendanceInfoResponse
	if personModel.TimeAttendanceID != nil && *personModel.TimeAttendanceID != "" {
		timeAttendanceModel, err := s.timeAttendanceRepo.GetByID(uuid.MustParse(*personModel.TimeAttendanceID))
		if err != nil {
			return nil, fmt.Errorf("failed to get time attendance: %w", err)
		}
		timeAttendance = &schema.AttendanceInfoResponse{
			ID:   timeAttendanceModel.ID.String(),
			Name: timeAttendanceModel.Name,
		}
	}

	cardIDs, err := s.personCardRepo.GetCardNumbersByPersonID(personModel.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get card numbers: %w", err)
	}

	licensePlateTexts, err := s.personLicenseRepo.GetLicensePlateTextsByPersonID(personModel.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get license plate texts: %w", err)
	}

	return &schema.PersonResponse{
		ID:                personModel.ID.String(),
		FirstName:         personModel.FirstName,
		MiddleName:        personModel.MiddleName,
		LastName:          personModel.LastName,
		PersonType:        personModel.PersonType,
		PersonID:          personModel.PersonID,
		Gender:            personModel.Gender,
		DateOfBirth:       personModel.DateOfBirth,
		Company:           personModel.Company,
		Department:        personModel.Department,
		JobPosition:       personModel.JobPosition,
		Address:           personModel.Address,
		MobileNumber:      personModel.MobileNumber,
		Email:             personModel.Email,
		IsVerified:        personModel.IsVerified,
		CardIDs:           cardIDs,
		LicensePlateTexts: licensePlateTexts,
		FaceImagePath:     personModel.FaceImagePath,
		ActiveAt:          personModel.ActiveAt,
		ExpireAt:          personModel.ExpireAt,
		AccessControlRule: accessRule,
		TimeAttendance:    timeAttendance,
	}, nil
}

// ----------> INNER FUNCTION <-----------------------//

func (s *personServiceImpl) validatePerson(isCreate bool, person *model.Person) error {
	if person.PersonID != nil && *person.PersonID != "" {
		isExist, err := s.personRepo.IsExistPersonID(*person.PersonID, person.ID)
		if err != nil {
			return fmt.Errorf("failed to check person ID existence: %w", err)
		}
		if isExist {
			return fmt.Errorf("person ID already exists")
		}
	}

	excludeID := uuid.Nil
	if !isCreate {
		excludeID = person.ID
	}
	isExistName, err := s.personRepo.IsExistName(person.FirstName, person.LastName, excludeID)
	if err != nil {
		return fmt.Errorf("failed to check person name existence: %w", err)
	}
	if isExistName {
		return fmt.Errorf("person with the same name already exists")
	}

	return nil
}

// handleRelationships manages the creation/deletion of related models within a transaction.
func (s *personServiceImpl) handleRelationships(
	txCardRepo repository.PersonCardRepository,
	txLicenseRepo repository.PersonLicensePlateRepository,
	personID string,
	cardIDs []string,
	licensePlateTexts []string,
) error {
	// 1. Delete old cards
	if err := txCardRepo.DeleteByPersonID(personID); err != nil {
		return fmt.Errorf("failed to delete old person cards: %w", err)
	}
	// 2. Create new cards
	if len(cardIDs) > 0 {
		cards := make([]model.PersonCard, len(cardIDs))
		for i, cardID := range cardIDs {
			cards[i] = model.PersonCard{CardNumber: cardID, PersonID: personID}
		}
		if err := txCardRepo.Create(cards); err != nil {
			return fmt.Errorf("failed to create new person cards: %w", err)
		}
	}

	// 3. Delete old license plates
	if err := txLicenseRepo.DeleteByPersonID(personID); err != nil {
		return fmt.Errorf("failed to delete old person license plates: %w", err)
	}
	// 4. Create new license plates
	if len(licensePlateTexts) > 0 {
		plates := make([]model.PersonLicensePlate, len(licensePlateTexts))
		for i, text := range licensePlateTexts {
			plates[i] = model.PersonLicensePlate{LicensePlateText: text, PersonID: personID}
		}
		if err := txLicenseRepo.Create(plates); err != nil {
			return fmt.Errorf("failed to create new person license plates: %w", err)
		}
	}

	return nil
}
