package admissionperiod

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/admissionperiod/handler"
	"go_be_enrollment/internal/modules/admissionperiod/repository"
	"go_be_enrollment/internal/modules/admissionperiod/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAdmissionPeriodRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewAdmissionPeriodRepository(db)
	svc := service.NewAdmissionPeriodService(repo)
	h := handler.NewAdmissionPeriodHandler(svc)

	// Admin API
	adminGroup := api.Group("/admin/admission-periods", middleware.AdminAuth(cfg))
	adminGroup.Get("/", h.GetList)
	adminGroup.Get("/:id", h.GetDetail)
	adminGroup.Post("/", h.Create)
	adminGroup.Put("/:id", h.Update)
	adminGroup.Patch("/:id/open-status", h.UpdateStatus)

	// Public API
	publicGroup := api.Group("/public/admission-periods")
	publicGroup.Get("/open", h.GetOpenPeriods)
}
