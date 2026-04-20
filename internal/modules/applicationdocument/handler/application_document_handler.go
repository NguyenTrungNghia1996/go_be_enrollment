package handler

import (
	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/applicationdocument/service"
	"go_be_enrollment/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ApplicationDocumentHandler struct {
	service service.ApplicationDocumentService
}

func NewApplicationDocumentHandler(service service.ApplicationDocumentService) *ApplicationDocumentHandler {
	return &ApplicationDocumentHandler{
		service: service,
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

func (h *ApplicationDocumentHandler) GetAdminList(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("id")
	if err != nil || appID <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID hồ sơ không hợp lệ", nil)
	}

	res, err := h.service.GetAdminList(uint(appID))
	if err != nil {
		logger.Log.Error("ApplicationDocumentHandler.GetAdminList", zap.Error(err))
		return httpresponse.ServerError(c)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách tài liệu thành công", res)
}

// ---------------- User Handlers ----------------

func (h *ApplicationDocumentHandler) GetUserList(c *fiber.Ctx) error {
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

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách tài liệu thành công", res)
}

func (h *ApplicationDocumentHandler) UploadDocument(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	appID, err := c.ParamsInt("id")
	if err != nil || appID <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID hồ sơ không hợp lệ", nil)
	}

	// Lấy thẻ Form value document_type
	docType := c.FormValue("document_type")
	if docType == "" {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng cung cấp loại tài liệu (document_type)", nil)
	}

	// Lấy file upload
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Vui lòng tải lên file đính kèm hợp lệ", nil)
	}

	res, err := h.service.UploadDocument(uint(appID), userID, docType, fileHeader)
	if err != nil {
		logger.Log.Error("ApplicationDocumentHandler.UploadDocument", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusCreated, "Tải tài liệu thành công", res)
}

func (h *ApplicationDocumentHandler) DeleteDocument(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Vui lòng đăng nhập lại", nil)
	}

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "ID tài liệu không hợp lệ", nil)
	}

	err = h.service.DeleteDocument(uint(id), userID)
	if err != nil {
		logger.Log.Error("ApplicationDocumentHandler.DeleteDocument", zap.Error(err))
		return httpresponse.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Đã xóa tài liệu thành công", nil)
}
