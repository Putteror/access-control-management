package repository

import (
	"github.com/putteror/access-control-management/internal/app/model"

	"gorm.io/gorm"
)

type PersonRepository interface {
	FindAll() ([]model.Person, error)
	FindByID(id string) (*model.Person, error)
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

func (r *personRepositoryImpl) FindByID(id string) (*model.Person, error) {
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
