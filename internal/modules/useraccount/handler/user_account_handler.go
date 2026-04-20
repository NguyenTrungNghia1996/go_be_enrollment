package handler

import (
	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/useraccount/dto"
	"go_be_enrollment/internal/modules/useraccount/service"
	"go_be_enrollment/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserAccountHandler struct {
	service service.UserAccountService
	val     *validator.Validate
}

func NewUserAccountHandler(service service.UserAccountService) *UserAccountHandler {
	return &UserAccountHandler{
		service: service,
		val:     validator.New(),
	}
}

func (h *UserAccountHandler) GetList(c *fiber.Ctx) error {
	filter := new(dto.UserAccountFilter)
	if err := c.QueryParser(filter); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Tham số không hợp lệ", nil)
	}

	res, err := h.service.GetList(filter)
	if err != nil {
		logger.Log.Error("UserAccountHandler.GetList", zap.Error(err))
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách thành công", res)
}

func (h *UserAccountHandler) GetDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID không hợp lệ", nil)
	}

	res, err := h.service.GetDetail(uint(id))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusNotFound, "NOT_FOUND", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy chi tiết thành công", res)
}

func (h *UserAccountHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID không hợp lệ", nil)
	}

	req := new(dto.UserAccountUpdateReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	if err := h.val.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng truyền đầy đủ tên", nil)
	}

	res, err := h.service.Update(uint(id), req)
	if err != nil {
		logger.Log.Error("UserAccountHandler.Update", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật tài khoản thành công", res)
}

func (h *UserAccountHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID không hợp lệ", nil)
	}

	req := new(dto.UserAccountStatusReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	err = h.service.UpdateStatus(uint(id), req)
	if err != nil {
		logger.Log.Error("UserAccountHandler.UpdateStatus", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật trạng thái thành công", nil)
}
