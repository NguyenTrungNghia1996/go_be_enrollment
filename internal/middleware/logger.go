package middleware

import (
	"time"

	"go_be_enrollment/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		cost := time.Since(start)

		logger.Log.Info("Incoming request",
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("query", string(c.Request().URI().QueryString())),
			zap.String("ip", c.IP()),
			zap.String("user-agent", string(c.Request().Header.UserAgent())),
			zap.Duration("cost", cost),
		)

		return err
	}
}
