package service

import (
	"errors"
	"math"
	"strings"

	"go_be_enrollment/internal/modules/useraccount/dto"
	"go_be_enrollment/internal/modules/useraccount/repository"

	authentity "go_be_enrollment/internal/modules/auth/entity"
)

type UserAccountService interface {
	GetList(filter *dto.UserAccountFilter) (*dto.PaginatedUserAccountRes, error)
	GetDetail(id uint) (*dto.UserAccountRes, error)
	Update(id uint, req *dto.UserAccountUpdateReq) (*dto.UserAccountRes, error)
	UpdateStatus(id uint, req *dto.UserAccountStatusReq) error
}

type userAccountService struct {
	repo repository.UserAccountRepository
}

func NewUserAccountService(repo repository.UserAccountRepository) UserAccountService {
	return &userAccountService{repo: repo}
}

func mapToDto(u *authentity.UserAccount) dto.UserAccountRes {
	return dto.UserAccountRes{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		FullName:      u.FullName,
		IsActive:      u.IsActive,
		TotalProfiles: 0, // TODO: Bổ sung count khi có bảng Profiles liên kết
		CreatedAt:     u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *userAccountService) GetList(filter *dto.UserAccountFilter) (*dto.PaginatedUserAccountRes, error) {
	list, total, err := s.repo.GetList(filter)
	if err != nil {
		return nil, errors.New("lỗi lấy danh sách tài khoản người dùng")
	}

	res := make([]dto.UserAccountRes, 0)
	for _, u := range list {
		res = append(res, mapToDto(&u))
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedUserAccountRes{
		Data:       res,
		Total:      total,
		Page:       filter.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *userAccountService) GetDetail(id uint) (*dto.UserAccountRes, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy tài khoản")
	}

	res := mapToDto(u)
	return &res, nil
}

func (s *userAccountService) Update(id uint, req *dto.UserAccountUpdateReq) (*dto.UserAccountRes, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy tài khoản")
	}

	u.FullName = strings.TrimSpace(req.FullName)
	u.Email = req.Email
	u.PhoneNumber = req.PhoneNumber

	// Clean emtpy pointers if any
	if u.Email != nil && strings.TrimSpace(*u.Email) == "" {
		u.Email = nil
	}
	if u.PhoneNumber != nil && strings.TrimSpace(*u.PhoneNumber) == "" {
		u.PhoneNumber = nil
	}

	if err := s.repo.Update(u); err != nil {
		return nil, errors.New("lỗi cập nhật tài khoản") // TODO: Bắt duplicate email/phone nếu cần
	}

	res := mapToDto(u)
	return &res, nil
}

func (s *userAccountService) UpdateStatus(id uint, req *dto.UserAccountStatusReq) error {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy tài khoản")
	}

	u.IsActive = req.IsActive

	if err := s.repo.Update(u); err != nil {
		return errors.New("lỗi cập nhật trạng thái tài khoản")
	}

	return nil
}
