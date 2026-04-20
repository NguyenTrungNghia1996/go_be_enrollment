package menu

import (
	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/middleware"
	
	adminauth_repo "go_be_enrollment/internal/modules/adminauth/repository"
	adminauth_service "go_be_enrollment/internal/modules/adminauth/service"
	
	"go_be_enrollment/internal/modules/menu/handler"
	"go_be_enrollment/internal/modules/menu/repository"
	"go_be_enrollment/internal/modules/menu/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterMenuRoutes(router fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := repository.NewMenuRepository(db)
	
	permRepo := adminauth_repo.NewPermissionRepository(db)
	permSvc := adminauth_service.NewPermissionService(permRepo)
	
	svc := service.NewMenuService(repo, permSvc)
	hdl := handler.NewMenuHandler(svc)

	menuGroup := router.Group("/admin/menus")

	menuGroup.Use(middleware.AdminAuth(cfg))
	
	// My Menu
	menuGroup.Get("/my-menu", hdl.GetMyMenu)
	
	// CRUD Menus
	menuGroup.Get("/", hdl.GetList)
	menuGroup.Get("/tree", hdl.GetTree)
	menuGroup.Post("/", hdl.Create)
	menuGroup.Put("/:id", hdl.Update)
	menuGroup.Delete("/:id", hdl.Delete)
}
