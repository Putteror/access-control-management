package main

import (
	"log"

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

	deviceRepo := repository.NewDeviceRepository(db)
	peopleRepo := repository.NewPeopleRepository(db)
	timerecordRepo := repository.NewTimeRecordRepository(db)

	deviceService := service.NewDeviceService(deviceRepo)
	peopleService := service.NewPeopleService(peopleRepo)
	timerecordService := service.NewTimeRecordService(timerecordRepo)

	deviceHandler := handler.NewDeviceHandler(deviceService)
	peopleHandler := handler.NewPeopleHandler(peopleService)
	timerecordHandler := handler.NewTimeRecordHandler(timerecordService)

	appRouter := router.NewRouter(deviceHandler, peopleHandler, timerecordHandler)

	log.Printf("Server is starting on port %s", cfg.Port)
	if err := appRouter.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Could not listen on %s: %v\n", cfg.Port, err)
	}
}
