package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

type AccessRecordRepository interface {
	GetAll(searchQuery schema.AccessRecordSearchQuery) ([]model.AccessRecord, error)
	GetByID(id uuid.UUID) (*model.AccessRecord, error)
	Create(accessRecord *model.AccessRecord) error
	Update(accessRecord *model.AccessRecord) error
	Delete(id uuid.UUID) error
}

type AccessRecordRepositoryImpl struct {
	db *gorm.DB
}

func NewAccessRecordRepository(db *gorm.DB) AccessRecordRepository {
	return &AccessRecordRepositoryImpl{db: db}
}

func (r *AccessRecordRepositoryImpl) GetAll(searchQuery schema.AccessRecordSearchQuery) ([]model.AccessRecord, error) {
	var accessRecords []model.AccessRecord

	query := r.db.Model(&model.AccessRecord{})

	var page int = searchQuery.Page
	var limit int = searchQuery.Limit
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&accessRecords).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated access records: %w", err)
	}

	return accessRecords, nil
}

func (r *AccessRecordRepositoryImpl) GetByID(id uuid.UUID) (*model.AccessRecord, error) {
	var accessRecord model.AccessRecord
	if err := r.db.First(&accessRecord, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &accessRecord, nil
}

func (r *AccessRecordRepositoryImpl) Create(accessRecord *model.AccessRecord) error {
	return r.db.Create(accessRecord).Error
}

func (r *AccessRecordRepositoryImpl) Update(accessRecord *model.AccessRecord) error {
	return r.db.Save(accessRecord).Error
}

func (r *AccessRecordRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Unscoped().Where("id = ?", id).Delete(&model.AccessRecord{}).Error
}
