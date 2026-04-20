package handler

import (
	"strconv"

	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/rolegroup/dto"
	"go_be_enrollment/internal/modules/rolegroup/service"

	"github.com/gofiber/fiber/v2"
)

type RoleGroupHandler struct {
	svc service.RoleGroupService
}

func NewRoleGroupHandler(svc service.RoleGroupService) *RoleGroupHandler {
	return &RoleGroupHandler{svc: svc}
}

func (h *RoleGroupHandler) GetList(c *fiber.Ctx) error {
	var filter dto.RoleGroupFilter
	if err := c.QueryParser(&filter); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_QUERY", "Tham số không hợp lệ", nil)
	}

	res, err := h.svc.GetList(&filter)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách thành công", res)
}

func (h *RoleGroupHandler) GetDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	res, err := h.svc.GetDetail(uint(id))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusNotFound, "NOT_FOUND", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy thông tin thành công", res)
}

func (h *RoleGroupHandler) Create(c *fiber.Ctx) error {
	var req dto.RoleGroupCreateReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	// TODO: Add Validator

	res, err := h.svc.Create(&req)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Tạo nhóm quyền thành công", res)
}

func (h *RoleGroupHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	var req dto.RoleGroupUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	// TODO: Add Validator

	res, err := h.svc.Update(uint(id), &req)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật thành công", res)
}

func (h *RoleGroupHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	var req dto.RoleGroupStatusReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	if err := h.svc.UpdateStatus(uint(id), &req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "UPDATE_STATUS_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật trạng thái thành công", nil)
}

func (h *RoleGroupHandler) GetPermissions(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	res, err := h.svc.GetPermissions(uint(id))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "GET_PERMISSIONS_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách permissions thành công", res)
}

func (h *RoleGroupHandler) UpdatePermissions(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	var req dto.RoleGroupPermissionReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	res, err := h.svc.ReplacePermissions(uint(id), &req)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "UPDATE_PERMISSIONS_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật permissions thành công", res)
}
