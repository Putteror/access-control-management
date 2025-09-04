package service

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
)

type PeopleService interface {
	CreatePeople(people *model.Person) error
	GetPeopleByID(id uint) (*model.Person, error)
	GetAllPeople() ([]model.Person, error)
}

type peopleServiceImpl struct {
	repo repository.PeopleRepository
}

func NewPeopleService(repo repository.PeopleRepository) PeopleService {
	return &peopleServiceImpl{repo: repo}
}

func (s *peopleServiceImpl) CreatePeople(people *model.Person) error {
	// Business Logic
	if people.FirstName == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return s.repo.Create(people)
}

func (s *peopleServiceImpl) GetPeopleByID(id uint) (*model.Person, error) {
	return s.repo.FindByID(id)
}

func (s *peopleServiceImpl) GetAllPeople() ([]model.Person, error) {
	return s.repo.FindAll()
}
