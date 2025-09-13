package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/putteror/access-control-management/internal/app/handler"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/service"
	"github.com/putteror/access-control-management/internal/config"
	"github.com/putteror/access-control-management/internal/database"
	"github.com/putteror/access-control-management/internal/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	database.AutoMigrate(db)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	uploadPath := filepath.Join(wd, "uploads")

	fileRepo := repository.NewFileSystemRepo(uploadPath)
	peopleRepo := repository.NewPersonRepository(db)
	accessControlRuleRepo := repository.NewAccessControlRuleRepository(db)

	peopleService := service.NewPersonService(peopleRepo, fileRepo, accessControlRuleRepo)

	peopleHandler := handler.NewPersonHandler(peopleService)

	appRouter := router.NewRouter(peopleHandler)

	log.Printf("Server is starting on port %s", cfg.Port)
	if err := appRouter.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Could not listen on %s: %v\n", cfg.Port, err)
	}
}
