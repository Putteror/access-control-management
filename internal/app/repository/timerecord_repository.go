package repository

import (
	"github.com/putteror/access-control-management/internal/app/model"

	"gorm.io/gorm"
)

type TimeRecordRepository interface {
	Create(timeRecord *model.TimeRecord) error
	FindAllByPeopleID(peopleID uint) ([]model.TimeRecord, error)
}

type timerecordRepositoryImpl struct {
	db *gorm.DB
}

func NewTimeRecordRepository(db *gorm.DB) TimeRecordRepository {
	return &timerecordRepositoryImpl{db: db}
}

func (r *timerecordRepositoryImpl) Create(timeRecord *model.TimeRecord) error {
	return r.db.Create(timeRecord).Error
}

func (r *timerecordRepositoryImpl) FindAllByPeopleID(peopleID uint) ([]model.TimeRecord, error) {
	var records []model.TimeRecord
	if err := r.db.Where("people_id = ?", peopleID).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
