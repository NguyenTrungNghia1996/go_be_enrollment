package middleware

import (
	"fmt"
	"runtime/debug"

	"go_be_enrollment/internal/common"
	"go_be_enrollment/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("%v", r)
				logger.Log.Error("Panic recovered",
					zap.Error(err),
					zap.String("stack", string(debug.Stack())),
				)
				_ = common.ErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error", nil)
			}
		}()
		return c.Next()
	}
}
