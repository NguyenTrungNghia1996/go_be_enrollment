package health

import (
	"go_be_enrollment/internal/common"

	"github.com/gofiber/fiber/v2"
)

func Check(c *fiber.Ctx) error {
	return common.SuccessResponse(c, fiber.StatusOK, "Service is up and running", fiber.Map{
		"status": "OK",
	})
}
