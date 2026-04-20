package handler

import (
	"strconv"

	"go_be_enrollment/internal/common/httpresponse"
	"go_be_enrollment/internal/modules/admissionperiodsubject/dto"
	"go_be_enrollment/internal/modules/admissionperiodsubject/service"

	"github.com/gofiber/fiber/v2"
)

type AdmissionPeriodSubjectHandler struct {
	svc service.AdmissionPeriodSubjectService
}

func NewAdmissionPeriodSubjectHandler(svc service.AdmissionPeriodSubjectService) *AdmissionPeriodSubjectHandler {
	return &AdmissionPeriodSubjectHandler{svc: svc}
}

func (h *AdmissionPeriodSubjectHandler) GetList(c *fiber.Ctx) error {
	periodID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID kỳ tuyển sinh không hợp lệ", nil)
	}

	res, err := h.svc.GetByAdmissionPeriodID(uint(periodID))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "GET_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Lấy danh sách môn học của kỳ tuyển sinh thành công", res)
}

func (h *AdmissionPeriodSubjectHandler) Replace(c *fiber.Ctx) error {
	periodID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_PARAM", "ID kỳ tuyển sinh không hợp lệ", nil)
	}

	var req dto.ReplaceAdmissionPeriodSubjectsReq
	if err := c.BodyParser(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "INVALID_BODY", "Dữ liệu không hợp lệ", nil)
	}

	res, err := h.svc.ReplaceSubjects(uint(periodID), &req)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "REPLACE_FAILED", err.Error(), nil)
	}

	return httpresponse.Success(c, fiber.StatusOK, "Cập nhật danh sách môn học cho kỳ tuyển sinh thành công", res)
}
