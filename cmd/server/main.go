package main

import (
	"fmt"
	"log"

	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/database"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/auth"
	"go_be_enrollment/internal/modules/health"
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

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Log.Info("Server is starting...", zap.String("port", cfg.AppPort))
	if err := app.Listen(addr); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
