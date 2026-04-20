package service

import (
	"errors"
	"math"
	"time"

	"go_be_enrollment/internal/modules/admissionperiod/dto"
	"go_be_enrollment/internal/modules/admissionperiod/entity"
	"go_be_enrollment/internal/modules/admissionperiod/repository"
)

type AdmissionPeriodService interface {
	GetList(filter *dto.AdmissionPeriodFilter) (*dto.PaginatedAdmissionPeriodRes, error)
	GetDetail(id uint) (*dto.AdmissionPeriodRes, error)
	Create(req *dto.AdmissionPeriodReq) (*dto.AdmissionPeriodRes, error)
	Update(id uint, req *dto.AdmissionPeriodReq) (*dto.AdmissionPeriodRes, error)
	UpdateStatus(id uint, req *dto.AdmissionPeriodStatusReq) error
	GetOpenPeriods() ([]dto.AdmissionPeriodRes, error)
}

type admissionPeriodService struct {
	repo repository.AdmissionPeriodRepository
}

func NewAdmissionPeriodService(repo repository.AdmissionPeriodRepository) AdmissionPeriodService {
	return &admissionPeriodService{repo: repo}
}

func parseDate(d *string) (*time.Time, error) {
	if d == nil || *d == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *d)
	if err != nil {
		return nil, errors.New("định dạng ngày không hợp lệ. Vui lòng dùng YYYY-MM-DD")
	}
	return &t, nil
}

func formatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func mapToDto(ap *entity.AdmissionPeriod) dto.AdmissionPeriodRes {
	return dto.AdmissionPeriodRes{
		ID:         ap.ID,
		Name:       ap.Name,
		SchoolYear: ap.SchoolYear,
		StartDate:  formatDate(ap.StartDate),
		EndDate:    formatDate(ap.EndDate),
		ExamFee:    ap.ExamFee,
		IsOpen:     ap.IsOpen,
		Notes:      ap.Notes,
		CreatedAt:  ap.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  ap.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *admissionPeriodService) GetList(filter *dto.AdmissionPeriodFilter) (*dto.PaginatedAdmissionPeriodRes, error) {
	list, total, err := s.repo.GetList(filter)
	if err != nil {
		return nil, errors.New("lỗi lấy danh sách kỳ tuyển sinh")
	}

	res := make([]dto.AdmissionPeriodRes, 0)
	for _, ap := range list {
		res = append(res, mapToDto(&ap))
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedAdmissionPeriodRes{
		Data:       res,
		Total:      total,
		Page:       filter.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *admissionPeriodService) GetDetail(id uint) (*dto.AdmissionPeriodRes, error) {
	ap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy kỳ tuyển sinh")
	}

	res := mapToDto(ap)
	return &res, nil
}

func (s *admissionPeriodService) validateBusinessRules(req *dto.AdmissionPeriodReq, start *time.Time, end *time.Time) error {
	if req.ExamFee < 0 {
		return errors.New("lệ phí thi không được nhỏ hơn 0")
	}

	if start != nil && end != nil {
		if end.Before(*start) {
			return errors.New("ngày kết thúc không được nhỏ hơn ngày bắt đầu")
		}
	}
	return nil
}

func (s *admissionPeriodService) Create(req *dto.AdmissionPeriodReq) (*dto.AdmissionPeriodRes, error) {
	start, err := parseDate(req.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := parseDate(req.EndDate)
	if err != nil {
		return nil, err
	}

	if err := s.validateBusinessRules(req, start, end); err != nil {
		return nil, err
	}

	ap := &entity.AdmissionPeriod{
		Name:       req.Name,
		SchoolYear: req.SchoolYear,
		StartDate:  start,
		EndDate:    end,
		ExamFee:    req.ExamFee,
		IsOpen:     req.IsOpen,
		Notes:      req.Notes,
	}

	if err := s.repo.Create(ap); err != nil {
		return nil, errors.New("lỗi tạo kỳ tuyển sinh")
	}

	res := mapToDto(ap)
	return &res, nil
}

func (s *admissionPeriodService) Update(id uint, req *dto.AdmissionPeriodReq) (*dto.AdmissionPeriodRes, error) {
	ap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy kỳ tuyển sinh")
	}

	start, err := parseDate(req.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := parseDate(req.EndDate)
	if err != nil {
		return nil, err
	}

	if err := s.validateBusinessRules(req, start, end); err != nil {
		return nil, err
	}

	ap.Name = req.Name
	ap.SchoolYear = req.SchoolYear
	ap.StartDate = start
	ap.EndDate = end
	ap.ExamFee = req.ExamFee
	ap.IsOpen = req.IsOpen
	ap.Notes = req.Notes

	if err := s.repo.Update(ap); err != nil {
		return nil, errors.New("lỗi cập nhật kỳ tuyển sinh")
	}

	res := mapToDto(ap)
	return &res, nil
}

func (s *admissionPeriodService) UpdateStatus(id uint, req *dto.AdmissionPeriodStatusReq) error {
	ap, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy kỳ tuyển sinh")
	}

	ap.IsOpen = req.IsOpen

	if err := s.repo.Update(ap); err != nil {
		return errors.New("lỗi cập nhật trạng thái kỳ tuyển sinh")
	}

	return nil
}

func (s *admissionPeriodService) GetOpenPeriods() ([]dto.AdmissionPeriodRes, error) {
	list, err := s.repo.GetOpenPeriods()
	if err != nil {
		return nil, errors.New("lỗi lấy danh sách kỳ tuyển sinh")
	}

	res := make([]dto.AdmissionPeriodRes, 0)
	for _, ap := range list {
		res = append(res, mapToDto(&ap))
	}
	return res, nil
}
