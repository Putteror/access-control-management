package repository

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"

	"gorm.io/gorm"
)

type PersonRepository interface {
	FindAll() ([]model.Person, error)
	FindByID(id string) (*model.Person, error)
	Create(people *model.Person) error
	Update(people *model.Person) error
	Delete(id string) error
	SaveImage(faceImageFile *multipart.FileHeader) (string, error)
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
	return r.db.Where("id = ?", id).Delete(&model.Person{}).Error
}

func (r *personRepositoryImpl) SaveImage(faceImageFile *multipart.FileHeader) (string, error) {

	uploadPath := "uploads/images/faces/people"
	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create a unique filename.
	filename := uuid.New().String() + filepath.Ext(faceImageFile.Filename)
	filePath := filepath.Join(uploadPath, filename)

	// Open the destination file.
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Open the source file.
	src, err := faceImageFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file header: %w", err)
	}
	defer src.Close()

	// Copy the file content.
	if _, err := fmt.Fprintf(dst, "%s", src); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	return filePath, nil
}
