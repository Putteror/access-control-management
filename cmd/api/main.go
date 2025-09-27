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

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("Error getting working directory: %v", err)
	// }
	// uploadPath := filepath.Join(wd, common.UploadPath)
	// fileRepo := repository.NewFileSystemRepo(uploadPath)

	accessControlDeviceRepo := repository.NewAccessControlDeviceRepository(db)
	accessControlGroupRepo := repository.NewAccessControlGroupRepository(db)
	accessControlRuleRepo := repository.NewAccessControlRuleRepository(db)
	accessControlServerRepo := repository.NewAccessControlServerRepository(db)
	accessRecordRepo := repository.NewAccessRecordRepository(db)
	AttendanceRepo := repository.NewAttendanceRepository(db)
	personRepo := repository.NewPersonRepository(db)
	userRepository := repository.NewUserRepository(db)

	accessControlDeviceService := service.NewAccessControlDeviceService(accessControlDeviceRepo, accessControlServerRepo)
	accessControlGroupService := service.NewAccessControlGroupService(accessControlGroupRepo, accessControlDeviceRepo, db)
	accessControlRuleService := service.NewAccessControlRuleService(accessControlRuleRepo, accessControlGroupRepo, db)
	accessRecordService := service.NewAccessRecordService(accessRecordRepo, personRepo, accessControlDeviceRepo)
	accessControlServerService := service.NewAccessControlServerService(accessControlServerRepo)
	attendanceService := service.NewAttendanceService(AttendanceRepo, db)
	authService := service.NewAuthService(userRepository)
	personService := service.NewPersonService(personRepo, repository.NewPersonCardRepository(db), repository.NewPersonLicensePlateRepository(db), accessControlRuleRepo, AttendanceRepo, db)
	userService := service.NewUserService(userRepository, db)

	accessControlDeviceHandler := handler.NewAccessControlDeviceHandler(accessControlDeviceService)
	accessControlGroupHandler := handler.NewAccessControlGroupHandler(accessControlGroupService)
	accessControlRuleHandler := handler.NewAccessControlRuleHandler(accessControlRuleService)
	accessControlServerHandler := handler.NewAccessControlServerHandler(accessControlServerService)
	accessRecordHandler := handler.NewAccessRecordHandler(accessRecordService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	authHandler := handler.NewAuthHandler(authService)
	personHandler := handler.NewPersonHandler(personService)
	userHandler := handler.NewUserHandler(userService)

	appRouter := router.NewRouter(
		accessControlDeviceHandler,
		accessControlGroupHandler,
		accessControlRuleHandler,
		accessControlServerHandler,
		accessRecordHandler,
		attendanceHandler,
		authHandler,
		personHandler,
		userHandler,
	)

	log.Printf("Server is starting on port %s", cfg.Port)
	if err := appRouter.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Could not listen on %s: %v\n", cfg.Port, err)
	}
}
