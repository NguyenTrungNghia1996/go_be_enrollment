package rolegroup

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	"go_be_enrollment/internal/modules/rolegroup/handler"
	"go_be_enrollment/internal/modules/rolegroup/repository"
	"go_be_enrollment/internal/modules/rolegroup/service"

	adminauth_repo "go_be_enrollment/internal/modules/adminauth/repository"
	adminauth_service "go_be_enrollment/internal/modules/adminauth/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoleGroupRoutes(router fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewRoleGroupRepository(db)
	svc := service.NewRoleGroupService(repo)
	hdl := handler.NewRoleGroupHandler(svc)

	permRepo := adminauth_repo.NewPermissionRepository(db)
	permSvc := adminauth_service.NewPermissionService(permRepo)

	roleGroup := router.Group("/admin/role-groups")

	roleGroup.Use(middleware.AdminAuth(cfg))
	// TODO: Thêm PermissionGuard vào các route của role-groups khi đã có key (Ví dụ: "role_group_menu")
	_ = permSvc 

	roleGroup.Get("/", hdl.GetList)
	roleGroup.Get("/:id", hdl.GetDetail)
	roleGroup.Post("/", hdl.Create)
	roleGroup.Put("/:id", hdl.Update)
	roleGroup.Patch("/:id/status", hdl.UpdateStatus)
}
