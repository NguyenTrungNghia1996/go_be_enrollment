package admissionperiodsubject

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/admissionperiodsubject/handler"
	"go_be_enrollment/internal/modules/admissionperiodsubject/repository"
	"go_be_enrollment/internal/modules/admissionperiodsubject/service"
	period_repo "go_be_enrollment/internal/modules/admissionperiod/repository"
	subject_repo "go_be_enrollment/internal/modules/subject/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAdmissionPeriodSubjectRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewAdmissionPeriodSubjectRepository(db)
	pRepo := period_repo.NewAdmissionPeriodRepository(db)
	sRepo := subject_repo.NewSubjectRepository(db)
	svc := service.NewAdmissionPeriodSubjectService(repo, pRepo, sRepo)
	h := handler.NewAdmissionPeriodSubjectHandler(svc)

	admin := api.Group("/admin/admission-periods/:id/subjects")
	admin.Use(middleware.AdminAuth(cfg))
	admin.Get("/", h.GetList)
	admin.Put("/", h.Replace)
}
