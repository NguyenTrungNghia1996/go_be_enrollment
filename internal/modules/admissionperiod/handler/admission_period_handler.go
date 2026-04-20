package handler

import (
	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/admissionperiod/dto"
	"go_be_enrollment/internal/modules/admissionperiod/service"
	"go_be_enrollment/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AdmissionPeriodHandler struct {
	service service.AdmissionPeriodService
	val     *validator.Validate
}

func NewAdmissionPeriodHandler(service service.AdmissionPeriodService) *AdmissionPeriodHandler {
	return &AdmissionPeriodHandler{
		service: service,
		val:     validator.New(),
	}
}

func (h *AdmissionPeriodHandler) GetList(c *fiber.Ctx) error {
	filter := new(dto.AdmissionPeriodFilter)
	if err := c.QueryParser(filter); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Tham số không hợp lệ", nil)
	}

	res, err := h.service.GetList(filter)
	if err != nil {
		logger.Log.Error("AdmissionPeriodHandler.GetList", zap.Error(err))
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách thành công", res)
}

func (h *AdmissionPeriodHandler) GetDetail(c *fiber.Ctx) error {
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

func (h *AdmissionPeriodHandler) Create(c *fiber.Ctx) error {
	req := new(dto.AdmissionPeriodReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	if err := h.val.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng truyền đầy đủ tên, năm học", nil)
	}

	res, err := h.service.Create(req)
	if err != nil {
		logger.Log.Error("AdmissionPeriodHandler.Create", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusCreated, "Tạo kỳ tuyển sinh thành công", res)
}

func (h *AdmissionPeriodHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID không hợp lệ", nil)
	}

	req := new(dto.AdmissionPeriodReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	if err := h.val.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng truyền đầy đủ tên, năm học", nil)
	}

	res, err := h.service.Update(uint(id), req)
	if err != nil {
		logger.Log.Error("AdmissionPeriodHandler.Update", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật kỳ tuyển sinh thành công", res)
}

func (h *AdmissionPeriodHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID không hợp lệ", nil)
	}

	req := new(dto.AdmissionPeriodStatusReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	err = h.service.UpdateStatus(uint(id), req)
	if err != nil {
		logger.Log.Error("AdmissionPeriodHandler.UpdateStatus", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật trạng thái thành công", nil)
}

func (h *AdmissionPeriodHandler) GetOpenPeriods(c *fiber.Ctx) error {
	res, err := h.service.GetOpenPeriods()
	if err != nil {
		logger.Log.Error("AdmissionPeriodHandler.GetOpenPeriods", zap.Error(err))
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách kỳ tuyển sinh đang mở thành công", res)
}
