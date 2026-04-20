package main

import (
	"fmt"
	"log"

	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/database"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/auth"
	"go_be_enrollment/internal/modules/auth/entity"
	"go_be_enrollment/internal/modules/adminauth"
	adminentity "go_be_enrollment/internal/modules/adminauth/entity"
	"go_be_enrollment/internal/modules/adminuser"
	"go_be_enrollment/internal/modules/rolegroup"
	"go_be_enrollment/internal/modules/menu"
	"go_be_enrollment/internal/modules/province"
	province_entity "go_be_enrollment/internal/modules/province/entity"
	"go_be_enrollment/internal/modules/wardunit"
	wardunit_entity "go_be_enrollment/internal/modules/wardunit/entity"
	"go_be_enrollment/internal/modules/admissionperiod"
	admission_period_entity "go_be_enrollment/internal/modules/admissionperiod/entity"
	"go_be_enrollment/internal/modules/application"
	application_entity "go_be_enrollment/internal/modules/application/entity"
	"go_be_enrollment/internal/modules/academicrecord"
	academic_record_entity "go_be_enrollment/internal/modules/academicrecord/entity"
	"go_be_enrollment/internal/modules/applicationdocument"
	application_document_entity "go_be_enrollment/internal/modules/applicationdocument/entity"
	"go_be_enrollment/internal/modules/health"
	"go_be_enrollment/internal/modules/useraccount"
	"go_be_enrollment/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to start application due to missing config: %v", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.AppEnv)
	defer logger.Sync()

	// Initialize database
	db, err := database.ConnectMySQL(cfg)
	if err != nil {
		logger.Log.Fatal("Could not initialize database connection", zap.Error(err))
	}
	defer database.Close()
	
	// AutoMigrate cho UserAccount
	if err := db.AutoMigrate(&entity.UserAccount{}); err != nil {
		logger.Log.Fatal("AutoMigrate failed for UserAccount", zap.Error(err))
	}
	
	// AutoMigrate cho AdminUser và quyền hạn
	if err := db.AutoMigrate(
		&adminentity.AdminUser{},
		&adminentity.RoleGroup{},
		&adminentity.AdminUserRoleGroup{},
		&adminentity.RoleGroupPermission{},
		&adminentity.Menu{},
		&province_entity.Province{},
		&wardunit_entity.WardUnit{},
		&admission_period_entity.AdmissionPeriod{},
		&application_entity.Application{},
		&academic_record_entity.AcademicRecord{},
		&application_document_entity.ApplicationDocument{},
	); err != nil {
		logger.Log.Fatal("AutoMigrate failed for Admin modules", zap.Error(err))
	}
	_ = db

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	// Apply Middlewares globally (Order matters!)
	app.Use(middleware.RequestID())       // Generate request ID first
	app.Use(middleware.RequestMetadata()) // Append contextual meta second
	app.Use(middleware.RequestLogger())   // Zap logger uses ID payload securely
	app.Use(middleware.Recovery())        // Catch Panics internally, outputting standard ID json
	app.Use(middleware.CORS(cfg))         // CORS dynamic configuration bound

	// Register Routes
	api := app.Group("/api/v1")
	health.RegisterRoutes(api)
	auth.RegisterUserAuthRoutes(api, db, cfg)
	adminauth.RegisterAdminAuthRoutes(api, db, cfg)
	adminuser.RegisterAdminUserRoutes(api, db, cfg)
	rolegroup.RegisterRoleGroupRoutes(api, db, cfg)
	menu.RegisterMenuRoutes(api, db, cfg)
	province.RegisterProvinceRoutes(api, db, cfg)
	wardunit.RegisterWardUnitRoutes(api, db, cfg)
	admissionperiod.RegisterAdmissionPeriodRoutes(api, db, cfg)
	useraccount.RegisterUserAccountRoutes(api, db, cfg)
	application.RegisterApplicationRoutes(api, db, cfg)
	academicrecord.RegisterAcademicRecordRoutes(api, db, cfg)
	applicationdocument.RegisterApplicationDocumentRoutes(api, db, cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Log.Info("Server is starting...", zap.String("port", cfg.AppPort))
	if err := app.Listen(addr); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
