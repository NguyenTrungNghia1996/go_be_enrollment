package academicrecord

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/academicrecord/handler"
	"go_be_enrollment/internal/modules/academicrecord/repository"
	"go_be_enrollment/internal/modules/academicrecord/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAcademicRecordRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewAcademicRecordRepository(db)
	svc := service.NewAcademicRecordService(repo)
	h := handler.NewAcademicRecordHandler(svc)

	adminGroup := api.Group("/admin/applications", middleware.AdminAuth(cfg))
	adminGroup.Get("/:id/academic-records", h.GetAdminList)

	userAppGroup := api.Group("/me/applications", middleware.UserAuth(cfg))
	userAppGroup.Get("/:id/academic-records", h.GetUserList)
	userAppGroup.Post("/:id/academic-records", h.Create)

	userRecGroup := api.Group("/me/academic-records", middleware.UserAuth(cfg))
	userRecGroup.Put("/:id", h.Update)
	userRecGroup.Delete("/:id", h.Delete)
}
