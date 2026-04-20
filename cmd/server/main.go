package main

import (
	"fmt"

	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/database"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/health"
	"go_be_enrollment/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger.InitLogger(cfg.AppEnv)
	defer logger.Sync()

	// Initialize database
	db, err := database.ConnectMySQL(cfg)
	if err != nil {
		logger.Log.Fatal("Could not initialize database connection", zap.Error(err))
	}
	// Prevent "declared and not used" error for now. `db` will be injected into repositories.
	_ = db

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Enrollment System",
	})

	// Apply Middlewares
	app.Use(middleware.RequestLogger())
	app.Use(middleware.Recovery())
	app.Use(middleware.CORS())

	// Register Routes
	api := app.Group("/api/v1")
	health.RegisterRoutes(api)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Log.Info(fmt.Sprintf("Server is starting on port %s...", cfg.AppPort))
	if err := app.Listen(addr); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
