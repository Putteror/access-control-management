package service

import (
	"fmt"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
)

type PeopleService interface {
	CreatePeople(people *model.People) error
	GetPeopleByID(id uint) (*model.People, error)
	GetAllPeople() ([]model.People, error)
}

type peopleServiceImpl struct {
	repo repository.PeopleRepository
}

func NewPeopleService(repo repository.PeopleRepository) PeopleService {
	return &peopleServiceImpl{repo: repo}
}

func (s *peopleServiceImpl) CreatePeople(people *model.People) error {
	// Business Logic
	if people.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return s.repo.Create(people)
}

func (s *peopleServiceImpl) GetPeopleByID(id uint) (*model.People, error) {
	return s.repo.FindByID(id)
}

func (s *peopleServiceImpl) GetAllPeople() ([]model.People, error) {
	return s.repo.FindAll()
}
