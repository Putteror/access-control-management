package service

import (
	"fmt"
	"mime/multipart"

	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

type PersonService interface {
	GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error)
	GetByID(id string) (*model.Person, error)
	Save(id string, personModel *model.Person, faceImageFile *multipart.FileHeader) error
	Delete(id string) error
	ConvertToResponse(person *model.Person) (*schema.PersonResponse, error)
}

type personServiceImpl struct {
	personRepo            repository.PersonRepository
	fileRepo              repository.FileRepository
	accessControlRuleRepo repository.AccessControlRuleRepository
}

func NewPersonService(personRepo repository.PersonRepository, fileRepo repository.FileRepository, accessControlRuleRepo repository.AccessControlRuleRepository) PersonService {
	return &personServiceImpl{
		personRepo:            personRepo,
		fileRepo:              fileRepo,
		accessControlRuleRepo: accessControlRuleRepo,
	}
}

func (s *personServiceImpl) GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error) {
	return s.personRepo.GetAll(searchQuery)
}

func (s *personServiceImpl) GetByID(id string) (*model.Person, error) {
	return s.personRepo.GetByID(id)
}

func (s *personServiceImpl) Save(id string, personModel *model.Person, faceImageFile *multipart.FileHeader) error {

	var faceImagePath string

	// validate
	if personModel.FirstName == "" || personModel.LastName == "" {
		return fmt.Errorf("firstName and lastName cannot be empty")
	}
	if faceImageFile != nil {
		var err error
		faceImagePath, err = s.fileRepo.Save(faceImageFile, common.FaceImagePath)
		if err != nil {
			return fmt.Errorf("failed to save image file: %w", err)
		}
		personModel.FaceImagePath = &faceImagePath
	}
	if personModel.AccessControlRuleID != nil && *personModel.AccessControlRuleID != "" {
		_, err := s.accessControlRuleRepo.GetByID(*personModel.AccessControlRuleID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("rule with ID '%s' does not exist", *personModel.AccessControlRuleID)
			}
			return fmt.Errorf("failed to retrieve access control rule: %w", err)
		}
	}

	// Save data
	if id == "" {
		err := s.personRepo.Create(personModel)
		if err != nil {
			if faceImagePath != "" {
				s.fileRepo.Delete(faceImagePath)
			}
			return fmt.Errorf("failed to save person to database: %w", err)
		}
	} else {
		err := s.personRepo.Update(personModel)
		if err != nil {
			if faceImagePath != "" {
				s.fileRepo.Delete(faceImagePath)
			}
			return fmt.Errorf("failed to update person in database: %w", err)
		}
	}

	return nil
}

func (s *personServiceImpl) Delete(id string) error {
	var faceImagePath string
	person, err := s.personRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("person with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get person by ID: %w", err)
	}
	if person.FaceImagePath != nil {
		faceImagePath = *person.FaceImagePath
		if err := s.fileRepo.Delete(faceImagePath); err != nil {
			return err
		}
	}
	return s.personRepo.Delete(id)
}

func (s *personServiceImpl) ConvertToResponse(person *model.Person) (*schema.PersonResponse, error) {
	var accessControlRuleResponse *schema.AccessControlRuleInfoResponse

	if person.AccessControlRuleID != nil && *person.AccessControlRuleID != "" {
		rule, err := s.accessControlRuleRepo.GetByID(*person.AccessControlRuleID)
		if err != nil {
			// If the rule is not found, we just return nil for the rule response without an error.
			if err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("failed to find access control rule with ID '%s': %w", *person.AccessControlRuleID, err)
			}
		}
		if rule != nil {
			accessControlRuleResponse = &schema.AccessControlRuleInfoResponse{
				ID:   rule.ID,
				Name: rule.Name,
			}
		}
	}

	response := &schema.PersonResponse{
		ID:                person.ID,
		FirstName:         person.FirstName,
		MiddleName:        person.MiddleName,
		LastName:          person.LastName,
		PersonType:        person.PersonType,
		PersonID:          person.PersonID,
		Gender:            person.Gender,
		DateOfBirth:       person.DateOfBirth,
		Company:           person.Company,
		Department:        person.Department,
		JobPosition:       person.JobPosition,
		Address:           person.Address,
		MobileNumber:      person.MobileNumber,
		Email:             person.Email,
		FaceImagePath:     person.FaceImagePath,
		ActiveAt:          person.ActiveAt,
		ExpireAt:          person.ExpireAt,
		AccessControlRule: accessControlRuleResponse,
	}

	return response, nil
}
