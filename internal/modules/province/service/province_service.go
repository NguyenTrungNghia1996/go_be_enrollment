package service

import (
	"errors"
	"math"

	"go_be_enrollment/internal/modules/province/dto"
	"go_be_enrollment/internal/modules/province/entity"
	"go_be_enrollment/internal/modules/province/repository"
)

type ProvinceService interface {
	GetList(filter *dto.ProvinceFilter) (*dto.PaginatedProvinceRes, error)
	GetDetail(id uint) (*dto.ProvinceRes, error)
	Create(req *dto.ProvinceCreateReq) (*dto.ProvinceRes, error)
	Update(id uint, req *dto.ProvinceUpdateReq) (*dto.ProvinceRes, error)
	UpdateStatus(id uint, req *dto.ProvinceStatusReq) error
	GetActiveList() ([]dto.ProvinceRes, error)
}

type provinceService struct {
	repo repository.ProvinceRepository
}

func NewProvinceService(repo repository.ProvinceRepository) ProvinceService {
	return &provinceService{repo: repo}
}

func mapToDto(p *entity.Province) dto.ProvinceRes {
	return dto.ProvinceRes{
		ID:        p.ID,
		Code:      p.Code,
		Name:      p.Name,
		IsActive:  p.IsActive,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *provinceService) GetList(filter *dto.ProvinceFilter) (*dto.PaginatedProvinceRes, error) {
	provinces, total, err := s.repo.GetList(filter)
	if err != nil {
		return nil, errors.New("lỗi lấy danh sách tỉnh/thành")
	}

	res := make([]dto.ProvinceRes, 0)
	for _, p := range provinces {
		res = append(res, mapToDto(&p))
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedProvinceRes{
		Data:       res,
		Total:      total,
		Page:       filter.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *provinceService) GetDetail(id uint) (*dto.ProvinceRes, error) {
	province, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy tỉnh/thành")
	}

	res := mapToDto(province)
	return &res, nil
}

func (s *provinceService) Create(req *dto.ProvinceCreateReq) (*dto.ProvinceRes, error) {
	if s.repo.CheckCodeExists(req.Code, 0) {
		return nil, errors.New("mã tỉnh/thành đã tồn tại")
	}
	if s.repo.CheckNameExists(req.Name, 0) {
		return nil, errors.New("tên tỉnh/thành đã tồn tại")
	}

	p := &entity.Province{
		Code:     req.Code,
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	if err := s.repo.Create(p); err != nil {
		return nil, errors.New("lỗi tạo tỉnh/thành")
	}

	res := mapToDto(p)
	return &res, nil
}

func (s *provinceService) Update(id uint, req *dto.ProvinceUpdateReq) (*dto.ProvinceRes, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy tỉnh/thành")
	}

	if s.repo.CheckCodeExists(req.Code, id) {
		return nil, errors.New("mã tỉnh/thành đã tồn tại")
	}
	if s.repo.CheckNameExists(req.Name, id) {
		return nil, errors.New("tên tỉnh/thành đã tồn tại")
	}

	p.Code = req.Code
	p.Name = req.Name
	p.IsActive = req.IsActive

	if err := s.repo.Update(p); err != nil {
		return nil, errors.New("lỗi cập nhật tỉnh/thành")
	}

	res := mapToDto(p)
	return &res, nil
}

func (s *provinceService) UpdateStatus(id uint, req *dto.ProvinceStatusReq) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy tỉnh/thành")
	}

	p.IsActive = req.IsActive

	if err := s.repo.Update(p); err != nil {
		return errors.New("lỗi cập nhật trạng thái tỉnh/thành")
	}

	return nil
}

func (s *provinceService) GetActiveList() ([]dto.ProvinceRes, error) {
	provinces, err := s.repo.GetActiveList()
	if err != nil {
		return nil, errors.New("lỗi lấy danh sách tỉnh/thành")
	}

	res := make([]dto.ProvinceRes, 0)
	for _, p := range provinces {
		res = append(res, mapToDto(&p))
	}
	return res, nil
}
