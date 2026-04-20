package handler

import (
	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/application/dto"
	"go_be_enrollment/internal/modules/application/service"
	"go_be_enrollment/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ApplicationHandler struct {
	service service.ApplicationService
	val     *validator.Validate
}

func NewApplicationHandler(service service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		service: service,
		val:     validator.New(),
	}
}

// ---------------- Admin Handlers ----------------

func (h *ApplicationHandler) GetAdminList(c *fiber.Ctx) error {
	filter := new(dto.ApplicationAdminFilter)
	if err := c.QueryParser(filter); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Tham số không hợp lệ", nil)
	}

	res, err := h.service.GetAdminList(filter)
	if err != nil {
		logger.Log.Error("ApplicationHandler.GetAdminList", zap.Error(err))
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách hồ sơ thành công", res)
}

func (h *ApplicationHandler) GetAdminDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID không hợp lệ", nil)
	}

	res, err := h.service.GetAdminDetail(uint(id))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusNotFound, "NOT_FOUND", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy chi tiết thành công", res)
}

// ---------------- User Handlers ----------------

func getUserID(c *fiber.Ctx) uint {
	v := c.Locals("user_id")
	switch val := v.(type) {
	case float64:
		return uint(val)
	case uint:
		return val
	case int:
		return uint(val)
	}
	return 0
}

func (h *ApplicationHandler) GetUserList(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	res, err := h.service.GetUserList(userID, page, limit)
	if err != nil {
		logger.Log.Error("ApplicationHandler.GetUserList", zap.Error(err))
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách hồ sơ cá nhân thành công", res)
}

func (h *ApplicationHandler) GetUserDetail(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID không hợp lệ", nil)
	}

	res, err := h.service.GetUserDetail(uint(id), userID)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusNotFound, "NOT_FOUND", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy chi tiết thành công", res)
}

func (h *ApplicationHandler) Create(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	req := new(dto.ApplicationReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	if err := h.val.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng điền đầy đủ các thông tin bắt buộc (Họ tên, Ngày sinh, CCCD...)", nil)
	}

	res, err := h.service.Create(userID, req)
	if err != nil {
		logger.Log.Error("ApplicationHandler.Create", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusCreated, "Tạo hồ sơ nháp thành công", res)
}

func (h *ApplicationHandler) Update(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID không hợp lệ", nil)
	}

	req := new(dto.ApplicationReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	if err := h.val.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng điền đầy đủ thông tin bắt buộc", nil)
	}

	res, err := h.service.Update(uint(id), userID, req)
	if err != nil {
		logger.Log.Error("ApplicationHandler.Update", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật hồ sơ thành công", res)
}
