package auth

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/auth/handler"
	"go_be_enrollment/internal/modules/auth/repository"
	"go_be_enrollment/internal/modules/auth/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserAuthRoutes(router fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewUserAccountRepository(db)
	svc := service.NewUserAuthService(repo, cfg)
	hdl := handler.NewUserAuthHandler(svc)

	authGroup := router.Group("/auth")
	{
		authGroup.Post("/register", hdl.Register)
		authGroup.Post("/login", hdl.Login)
		authGroup.Post("/activate", hdl.Activate)
		authGroup.Get("/activate", hdl.Activate)

		// Protected endpoints passing through User JWT Middleware
		protected := authGroup.Group("")
		protected.Use(middleware.UserAuth(cfg))
		protected.Get("/me", hdl.GetMe)
	}
}
