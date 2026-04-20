package middleware

import (
	"strings"

	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/config"
	"go_be_enrollment/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// UserAuth parses and validates the JWT Token for end-users
func UserAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid token", nil)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseUserToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			return httpresponse.Error(c, fiber.StatusUnauthorized, "TOKEN_EXPIRED_OR_INVALID", "Token expired or invalid", nil)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
