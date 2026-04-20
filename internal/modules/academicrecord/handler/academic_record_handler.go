package handler

import (
	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/academicrecord/dto"
	"go_be_enrollment/internal/modules/academicrecord/service"
	"go_be_enrollment/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AcademicRecordHandler struct {
	service service.AcademicRecordService
	val     *validator.Validate
}

func NewAcademicRecordHandler(service service.AcademicRecordService) *AcademicRecordHandler {
	return &AcademicRecordHandler{
		service: service,
		val:     validator.New(),
	}
}

// ---------------- Helper ----------------

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

// ---------------- Admin Handlers ----------------

func (h *AcademicRecordHandler) GetAdminList(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil || appID <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID hồ sơ không hợp lệ", nil)
	}

	res, err := h.service.GetAdminList(uint(appID))
	if err != nil {
		logger.Log.Error("AcademicRecordHandler.GetAdminList", zap.Error(err))
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách học bạ thành công", res)
}

// ---------------- User Handlers ----------------

func (h *AcademicRecordHandler) GetUserList(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	appID, err := c.ParamsInt("id")
	if err != nil || appID <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID hồ sơ không hợp lệ", nil)
	}

	res, err := h.service.GetUserList(uint(appID), userID)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách học bạ thành công", res)
}

func (h *AcademicRecordHandler) Create(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	appID, err := c.ParamsInt("id")
	if err != nil || appID <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID hồ sơ không hợp lệ", nil)
	}

	req := new(dto.AcademicRecordReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	if err := h.val.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng điền đầy đủ các thông tin bắt buộc (Lớp, Khối, Tên trường...)", nil)
	}

	res, err := h.service.Create(uint(appID), userID, req)
	if err != nil {
		logger.Log.Error("AcademicRecordHandler.Create", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusCreated, "Thêm học bạ thành công", res)
}

func (h *AcademicRecordHandler) Update(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID học bạ không hợp lệ", nil)
	}

	req := new(dto.AcademicRecordReq)
	if err := c.BodyParser(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Dữ liệu không hợp lệ", nil)
	}

	if err := h.val.Struct(req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng điền đầy đủ thông tin bắt buộc", nil)
	}

	res, err := h.service.Update(uint(id), userID, req)
	if err != nil {
		logger.Log.Error("AcademicRecordHandler.Update", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật học bạ thành công", res)
}

func (h *AcademicRecordHandler) Delete(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID học bạ không hợp lệ", nil)
	}

	err = h.service.Delete(uint(id), userID)
	if err != nil {
		logger.Log.Error("AcademicRecordHandler.Delete", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Đã xóa học bạ thành công", nil)
}
