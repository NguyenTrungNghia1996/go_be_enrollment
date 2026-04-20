package handler

import (
	"strconv"

	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/menu/dto"
	"go_be_enrollment/internal/modules/menu/service"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	svc service.MenuService
}

func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

func (h *MenuHandler) GetList(c *fiber.Ctx) error {
	res, err := h.svc.GetList()
	if err != nil {
		return httpresponse.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}
	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách menu thành công", res)
}

func (h *MenuHandler) GetTree(c *fiber.Ctx) error {
	res, err := h.svc.GetTree()
	if err != nil {
		return httpresponse.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}
	return httpresponse.Success(c, fiber.StatusOK, "Lấy cây menu thành công", res)
}

func (h *MenuHandler) GetMyMenu(c *fiber.Ctx) error {
	adminIDVal := c.Locals("admin_id")
	isSuperAdminVal := c.Locals("admin_is_super_admin")

	if adminIDVal == nil || isSuperAdminVal == nil {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Missing admin token", nil)
	}

	adminID := adminIDVal.(uint)
	isSuperAdmin := isSuperAdminVal.(bool)

	res, err := h.svc.GetMyMenuTree(adminID, isSuperAdmin)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}
	return httpresponse.Success(c, fiber.StatusOK, "Lấy cây menu của tôi thành công", res)
}

func (h *MenuHandler) Create(c *fiber.Ctx) error {
	var req dto.MenuCreateReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	res, err := h.svc.Create(&req)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Tạo menu thành công", res)
}

func (h *MenuHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	var req dto.MenuUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	res, err := h.svc.Update(uint(id), &req)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật menu thành công", res)
}

func (h *MenuHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Xóa menu thành công", nil)
}
