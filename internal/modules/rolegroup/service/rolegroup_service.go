package service

import (
	"errors"
	"math"

	"go_be_enrollment/internal/modules/adminauth/entity"
	"go_be_enrollment/internal/modules/rolegroup/dto"
	"go_be_enrollment/internal/modules/rolegroup/repository"
)

type RoleGroupService interface {
	GetList(filter *dto.RoleGroupFilter) (*dto.PaginatedRoleGroupRes, error)
	GetDetail(id uint) (*dto.RoleGroupRes, error)
	Create(req *dto.RoleGroupCreateReq) (*dto.RoleGroupRes, error)
	Update(id uint, req *dto.RoleGroupUpdateReq) (*dto.RoleGroupRes, error)
	UpdateStatus(id uint, req *dto.RoleGroupStatusReq) error
}

type roleGroupService struct {
	repo repository.RoleGroupRepository
}

func NewRoleGroupService(repo repository.RoleGroupRepository) RoleGroupService {
	return &roleGroupService{repo: repo}
}

func mapToDto(m *entity.RoleGroup) *dto.RoleGroupRes {
	return &dto.RoleGroupRes{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *roleGroupService) GetList(filter *dto.RoleGroupFilter) (*dto.PaginatedRoleGroupRes, error) {
	roles, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	var data []dto.RoleGroupRes
	for _, r := range roles {
		data = append(data, *mapToDto(&r))
	}
	if data == nil {
		data = []dto.RoleGroupRes{}
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 10
	}

	return &dto.PaginatedRoleGroupRes{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}, nil
}

func (s *roleGroupService) GetDetail(id uint) (*dto.RoleGroupRes, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy nhóm quyền")
	}
	return mapToDto(role), nil
}

func (s *roleGroupService) Create(req *dto.RoleGroupCreateReq) (*dto.RoleGroupRes, error) {
	if s.repo.CheckCodeExists(req.Code, 0) {
		return nil, errors.New("mã nhóm quyền đã tồn tại")
	}

	newRole := &entity.RoleGroup{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.repo.Create(newRole); err != nil {
		return nil, errors.New("lỗi tạo nhóm quyền")
	}

	return mapToDto(newRole), nil
}

func (s *roleGroupService) Update(id uint, req *dto.RoleGroupUpdateReq) (*dto.RoleGroupRes, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy nhóm quyền")
	}

	role.Name = req.Name
	role.Description = req.Description

	if err := s.repo.Update(role); err != nil {
		return nil, errors.New("lỗi cập nhật nhóm quyền")
	}

	return mapToDto(role), nil
}

func (s *roleGroupService) UpdateStatus(id uint, req *dto.RoleGroupStatusReq) error {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy nhóm quyền")
	}

	// Lock prevention if assigned
	if !req.IsActive && role.IsActive {
		if s.repo.CheckIsAssigned(id) {
			return errors.New("không thể khóa vì nhóm quyền này đang được gán cho quản trị viên")
		}
	}

	role.IsActive = req.IsActive
	if err := s.repo.Update(role); err != nil {
		return errors.New("lỗi cập nhật trạng thái nhóm quyền")
	}

	return nil
}
