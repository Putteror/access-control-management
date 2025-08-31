package database

import (
	"fmt"
	"log"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.People{}, &model.TimeRecord{})
}
