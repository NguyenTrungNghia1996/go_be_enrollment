package examroom

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/examroom/handler"
	"go_be_enrollment/internal/modules/examroom/repository"
	"go_be_enrollment/internal/modules/examroom/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterExamRoomRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewExamRoomRepository(db)
	svc := service.NewExamRoomService(repo)
	h := handler.NewExamRoomHandler(svc)

	admin := api.Group("/admin/exam-rooms")
	admin.Use(middleware.AdminAuth(cfg))
	admin.Get("/", h.GetList)
	admin.Get("/:id", h.GetDetail)
	admin.Post("/", h.Create)
	admin.Put("/:id", h.Update)
	admin.Delete("/:id", h.Delete)
}
