package health

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	healthGroup := router.Group("/health")
	{
		healthGroup.Get("", Check)
	}
}
