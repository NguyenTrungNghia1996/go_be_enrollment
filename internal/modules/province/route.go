package province

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/province/handler"
	"go_be_enrollment/internal/modules/province/repository"
	"go_be_enrollment/internal/modules/province/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterProvinceRoutes(router fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewProvinceRepository(db)
	svc := service.NewProvinceService(repo)
	hdl := handler.NewProvinceHandler(svc)

	// Admin routes
	adminGroup := router.Group("/admin/provinces")
	adminGroup.Use(middleware.AdminAuth(cfg))
	adminGroup.Get("/", hdl.GetList)
	adminGroup.Get("/:id", hdl.GetDetail)
	adminGroup.Post("/", hdl.Create)
	adminGroup.Put("/:id", hdl.Update)
	adminGroup.Patch("/:id/status", hdl.UpdateStatus)

	// Public routes
	publicGroup := router.Group("/public/provinces")
	publicGroup.Get("/", hdl.GetPublicList)
}
