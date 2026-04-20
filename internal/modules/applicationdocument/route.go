package applicationdocument

import (
	"go_be_enrollment/internal/common/storage"
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/applicationdocument/handler"
	"go_be_enrollment/internal/modules/applicationdocument/repository"
	"go_be_enrollment/internal/modules/applicationdocument/service"
	"go_be_enrollment/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterApplicationDocumentRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	storageSvc, err := storage.NewS3StorageService(cfg)
	if err != nil {
		logger.Log.Warn("Error initialize S3Storage", zap.Error(err))
	}

	repo := repository.NewApplicationDocumentRepository(db)
	svc := service.NewApplicationDocumentService(repo, storageSvc)
	h := handler.NewApplicationDocumentHandler(svc)

	adminGroup := api.Group("/admin/applications", middleware.AdminAuth(cfg))
	adminGroup.Get("/:id/documents", h.GetAdminList)

	userAppGroup := api.Group("/me/applications", middleware.UserAuth(cfg))
	userAppGroup.Get("/:id/documents", h.GetUserList)
	// Accept multipart/form-data
	userAppGroup.Post("/:id/documents", h.UploadDocument)

	userDocGroup := api.Group("/me/documents", middleware.UserAuth(cfg))
	userDocGroup.Delete("/:id", h.DeleteDocument)
}
