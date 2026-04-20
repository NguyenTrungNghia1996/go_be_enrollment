package handler

import (
	"strconv"

	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/examroom/dto"
	"go_be_enrollment/internal/modules/examroom/service"

	"github.com/gofiber/fiber/v2"
)

type ExamRoomHandler struct {
	svc service.ExamRoomService
}

func NewExamRoomHandler(svc service.ExamRoomService) *ExamRoomHandler {
	return &ExamRoomHandler{svc: svc}
}

func (h *ExamRoomHandler) GetList(c *fiber.Ctx) error {
	var filter dto.ExamRoomFilter
	if err := c.QueryParser(&filter); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_QUERY", "Tham số không hợp lệ", nil)
	}

	res, err := h.svc.GetList(&filter)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách phòng thi thành công", res)
}

func (h *ExamRoomHandler) GetDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	res, err := h.svc.GetDetail(uint(id))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusNotFound, "NOT_FOUND", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy thông tin phòng thi thành công", res)
}

func (h *ExamRoomHandler) Create(c *fiber.Ctx) error {
	var req dto.ExamRoomCreateReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	res, err := h.svc.Create(&req)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Tạo phòng thi thành công", res)
}

func (h *ExamRoomHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	var req dto.ExamRoomUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	res, err := h.svc.Update(uint(id), &req)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật phòng thi thành công", res)
}

func (h *ExamRoomHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID không hợp lệ", nil)
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Xóa phòng thi thành công", nil)
}
