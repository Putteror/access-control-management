package service

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
)

type PersonService interface {
	GetAllPeople() ([]model.Person, error)
	GetPersonByID(id string) (*model.Person, error)
	CreatePerson(person *model.Person, faceImageFile *multipart.FileHeader) error
	UpdatePerson(person *model.Person) error
	DeletePerson(id string) error
}

type personServiceImpl struct {
	repo repository.PersonRepository
}

func NewPersonService(repo repository.PersonRepository) PersonService {
	return &personServiceImpl{repo: repo}
}

func (s *personServiceImpl) GetAllPeople() ([]model.Person, error) {
	return s.repo.FindAll()
}

func (s *personServiceImpl) GetPersonByID(id string) (*model.Person, error) {
	return s.repo.FindByID(id)
}

func (s *personServiceImpl) CreatePerson(person *model.Person, faceImageFile *multipart.FileHeader) error {
	// Business Logic
	if person.FirstName == "" {
		return fmt.Errorf("name cannot be empty")
	}
	filePath, err := s.repo.SaveImage(faceImageFile)
	if err != nil {
		return fmt.Errorf("failed to save image file: %w", err)
	}

	// 2. Defer the file deletion in case of DB failure
	defer os.Remove(filePath)

	// 3. Update the person model
	person.FaceImagePath = filePath
	return s.repo.Create(person)
}

func (s *personServiceImpl) UpdatePerson(person *model.Person) error {
	return s.repo.Update(person)
}

func (s *personServiceImpl) DeletePerson(id string) error {
	return s.repo.Delete(id)
}
