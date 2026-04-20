package application

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/application/handler"
	"go_be_enrollment/internal/modules/application/repository"
	"go_be_enrollment/internal/modules/application/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterApplicationRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewApplicationRepository(db)
	svc := service.NewApplicationService(repo)
	h := handler.NewApplicationHandler(svc)

	adminGroup := api.Group("/admin/applications", middleware.AdminAuth(cfg))
	adminGroup.Get("/", h.GetAdminList)
	adminGroup.Get("/:id", h.GetAdminDetail)

	userGroup := api.Group("/me/applications", middleware.UserAuth(cfg))
	userGroup.Get("/", h.GetUserList)
	userGroup.Get("/:id", h.GetUserDetail)
	userGroup.Post("/", h.Create)
	userGroup.Put("/:id", h.Update)
}
