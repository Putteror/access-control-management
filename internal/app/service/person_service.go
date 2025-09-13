package service

import (
	"fmt"
	"mime/multipart"

	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
)

type PersonService interface {
	GetAll(searchQuery schema.PersonSearchQuery) ([]model.Person, error)
	GetByID(id string) (*model.Person, error)
	Save(id *string, personModel *model.Person, faceImageFile *multipart.FileHeader) error
	Delete(id string) error
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
	return s.personRepo.FindByID(id)
}

func (s *personServiceImpl) Save(id *string, personModel *model.Person, faceImageFile *multipart.FileHeader) error {

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
		personModel.FaceImagePath = faceImagePath
	}
	if personModel.AccessControlRuleID != "" {
		isExist, err := s.accessControlRuleRepo.IsExistID(personModel.AccessControlRuleID)
		if err != nil {
			return fmt.Errorf("failed to check if rule ID exists: %w", err)
		}
		if !isExist {
			return fmt.Errorf("rule ID '%s' does not exist", personModel.AccessControlRuleID)
		}
	}

	// Save data
	if id != nil {
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
	return s.personRepo.Delete(id)
}
