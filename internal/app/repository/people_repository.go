package repository

import (
	"github.com/putteror/access-control-management/internal/app/model"

	"gorm.io/gorm"
)

type PeopleRepository interface {
	Create(people *model.People) error
	FindByID(id uint) (*model.People, error)
	FindAll() ([]model.People, error)
	Update(people *model.People) error
	Delete(id uint) error
}

type peopleRepositoryImpl struct {
	db *gorm.DB
}

func NewPeopleRepository(db *gorm.DB) PeopleRepository {
	return &peopleRepositoryImpl{db: db}
}

func (r *peopleRepositoryImpl) Create(people *model.People) error {
	return r.db.Create(people).Error
}

func (r *peopleRepositoryImpl) FindByID(id uint) (*model.People, error) {
	var people model.People
	if err := r.db.First(&people, id).Error; err != nil {
		return nil, err
	}
	return &people, nil
}

func (r *peopleRepositoryImpl) FindAll() ([]model.People, error) {
	var peoples []model.People
	if err := r.db.Find(&peoples).Error; err != nil {
		return nil, err
	}
	return peoples, nil
}

func (r *peopleRepositoryImpl) Update(people *model.People) error {
	return r.db.Save(people).Error
}

func (r *peopleRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&model.People{}, id).Error
}
