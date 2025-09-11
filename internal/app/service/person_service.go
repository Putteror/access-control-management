package service

import (
	"fmt"
	"mime/multipart"

	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
)

type PersonService interface {
	GetAllPeople() ([]model.Person, error)
	PaginatedFindAllPeople(page, limit int) ([]model.Person, error)
	GetPersonByID(id string) (*model.Person, error)
	CreatePerson(person *model.Person, faceImageFile *multipart.FileHeader) error
	UpdatePerson(person *model.Person) error
	DeletePerson(id string) error
}

type personServiceImpl struct {
	personRepo repository.PersonRepository
	fileRepo   repository.FileRepository
}

func NewPersonService(personRepo repository.PersonRepository, fileRepo repository.FileRepository) PersonService {
	return &personServiceImpl{
		personRepo: personRepo,
		fileRepo:   fileRepo,
	}
}

func (s *personServiceImpl) GetAllPeople() ([]model.Person, error) {
	return s.personRepo.FindAll()
}

func (s *personServiceImpl) PaginatedFindAllPeople(page, limit int) ([]model.Person, error) {
	return s.personRepo.PaginatedFindAll(page, limit)
}

func (s *personServiceImpl) GetPersonByID(id string) (*model.Person, error) {
	return s.personRepo.FindByID(id)
}

func (s *personServiceImpl) CreatePerson(person *model.Person, faceImageFile *multipart.FileHeader) error {
	var filePath string
	if person.FirstName == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if faceImageFile != nil {
		var err error
		filePath, err = s.fileRepo.Save(faceImageFile, common.FaceImagePath)
		if err != nil {
			return fmt.Errorf("failed to save image file: %w", err)
		}
		person.FaceImagePath = filePath
	}

	err := s.personRepo.Create(person)
	if err != nil {
		if filePath != "" {
			s.fileRepo.Delete(filePath)
		}
		return fmt.Errorf("failed to save person to database: %w", err)
	}

	return nil
}

func (s *personServiceImpl) UpdatePerson(person *model.Person) error {
	return s.personRepo.Update(person)
}

func (s *personServiceImpl) DeletePerson(id string) error {
	return s.personRepo.Delete(id)
}
