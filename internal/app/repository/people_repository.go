package repository

import (
	"github.com/putteror/access-control-management/internal/app/model"

	"gorm.io/gorm"
)

type PeopleRepository interface {
	Create(people *model.Person) error
	FindByID(id uint) (*model.Person, error)
	FindAll() ([]model.Person, error)
	Update(people *model.Person) error
	Delete(id uint) error
}

type peopleRepositoryImpl struct {
	db *gorm.DB
}

func NewPeopleRepository(db *gorm.DB) PeopleRepository {
	return &peopleRepositoryImpl{db: db}
}

func (r *peopleRepositoryImpl) Create(people *model.Person) error {
	return r.db.Create(people).Error
}

func (r *peopleRepositoryImpl) FindByID(id uint) (*model.Person, error) {
	var people model.Person
	if err := r.db.First(&people, id).Error; err != nil {
		return nil, err
	}
	return &people, nil
}

func (r *peopleRepositoryImpl) FindAll() ([]model.Person, error) {
	var peoples []model.Person
	if err := r.db.Find(&peoples).Error; err != nil {
		return nil, err
	}
	return peoples, nil
}

func (r *peopleRepositoryImpl) Update(people *model.Person) error {
	return r.db.Save(people).Error
}

func (r *peopleRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&model.Person{}, id).Error
}
