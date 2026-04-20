package useraccount

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/useraccount/handler"
	"go_be_enrollment/internal/modules/useraccount/repository"
	"go_be_enrollment/internal/modules/useraccount/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserAccountRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewUserAccountRepository(db)
	svc := service.NewUserAccountService(repo)
	h := handler.NewUserAccountHandler(svc)

	adminGroup := api.Group("/admin/user-accounts", middleware.AdminAuth(cfg))
	adminGroup.Get("/", h.GetList)
	adminGroup.Get("/:id", h.GetDetail)
	adminGroup.Put("/:id", h.Update)
	adminGroup.Patch("/:id/status", h.UpdateStatus)
}
