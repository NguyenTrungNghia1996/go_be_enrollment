package handler

import (
	"errors"

	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/auth/dto"
	"go_be_enrollment/internal/modules/auth/service"
	"go_be_enrollment/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type UserAuthHandler struct {
	service service.UserAuthService
}

func NewUserAuthHandler(svc service.UserAuthService) *UserAuthHandler {
	return &UserAuthHandler{service: svc}
}

func (h *UserAuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Dữ liệu body không hợp lệ", nil)
	}

	if err := utils.Validate.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Vui lòng kiểm tra lại thông tin đầu vào", err.Error())
	}

	if err := h.service.Register(c.Context(), &req); err != nil {
		if errors.Is(err, service.ErrUsernameExists) || errors.Is(err, service.ErrEmailExists) {
			return httpresponse.Error(c, fiber.StatusConflict, "CONFLICT_DATA", err.Error(), nil)
		}
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusCreated, "Đăng ký thành công", nil)
}

func (h *UserAuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Dữ liệu body không hợp lệ", nil)
	}

	if err := utils.Validate.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Thiếu username hoặc password", nil)
	}

	resp, err := h.service.Login(c.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidLogin):
			return httpresponse.Error(c, fiber.StatusUnauthorized, "LOGIN_FAILED", "Tài khoản hoặc mật khẩu không đúng", nil)
		case errors.Is(err, service.ErrUserInactive):
			return httpresponse.Error(c, fiber.StatusForbidden, "ACCOUNT_LOCKED", "Tài khoản đã bị khóa", nil)
		default:
			return httpresponse.ServerError(c)
		}
	}

	return httpresponse.Success(c, fiber.StatusOK, "Đăng nhập thành công", resp)
}

func (h *UserAuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	resp, err := h.service.GetMe(c.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return httpresponse.Error(c, fiber.StatusNotFound, "NOT_FOUND", "Tài khoản không tồn tại", nil)
		}
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy thông tin thành công", resp)
}
