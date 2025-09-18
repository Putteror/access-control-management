package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/putteror/access-control-management/internal/app/common"
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
	uploadPath := filepath.Join(wd, common.UploadPath)

	fileRepo := repository.NewFileSystemRepo(uploadPath)
	accessControlDeviceRepo := repository.NewAccessControlDeviceRepository(db)
	accessControlGroupRepo := repository.NewAccessControlGroupRepository(db)
	accessControlRuleRepo := repository.NewAccessControlRuleRepository(db)
	accessControlServerRepo := repository.NewAccessControlServerRepository(db)
	AttendanceRepo := repository.NewAttendanceRepository(db)
	peopleRepo := repository.NewPersonRepository(db)

	accessControlDeviceService := service.NewAccessControlDeviceService(accessControlDeviceRepo, accessControlServerRepo)
	accessControlGroupService := service.NewAccessControlGroupService(accessControlGroupRepo, accessControlDeviceRepo, db)
	accessControlRuleService := service.NewAccessControlRuleService(accessControlRuleRepo, accessControlGroupRepo, db)
	accessControlServerService := service.NewAccessControlServerService(accessControlServerRepo)
	attendanceService := service.NewAttendanceService(AttendanceRepo, db)
	authService := service.NewAuthService()
	peopleService := service.NewPersonService(peopleRepo, fileRepo, accessControlRuleRepo)

	accessControlDeviceHandler := handler.NewAccessControlDeviceHandler(accessControlDeviceService)
	accessControlGroupHandler := handler.NewAccessControlGroupHandler(accessControlGroupService)
	accessControlRuleHandler := handler.NewAccessControlRuleHandler(accessControlRuleService)
	accessControlServerHandler := handler.NewAccessControlServerHandler(accessControlServerService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	authHandler := handler.NewAuthHandler(authService)
	peopleHandler := handler.NewPersonHandler(peopleService)

	appRouter := router.NewRouter(
		accessControlDeviceHandler,
		accessControlGroupHandler,
		accessControlRuleHandler,
		accessControlServerHandler,
		attendanceHandler,
		authHandler,
		peopleHandler,
	)

	log.Printf("Server is starting on port %s", cfg.Port)
	if err := appRouter.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Could not listen on %s: %v\n", cfg.Port, err)
	}
}
