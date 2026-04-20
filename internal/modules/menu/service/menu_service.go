package service

import (
	"errors"

	"go_be_enrollment/internal/modules/adminauth/entity"
	adminauth_service "go_be_enrollment/internal/modules/adminauth/service"
	"go_be_enrollment/internal/modules/menu/dto"
	"go_be_enrollment/internal/modules/menu/repository"
)

type MenuService interface {
	GetList() ([]dto.MenuRes, error)
	GetTree() ([]dto.MenuRes, error)
	GetMyMenuTree(adminID uint, isSuperAdmin bool) ([]dto.MenuRes, error)
	Create(req *dto.MenuCreateReq) (*dto.MenuRes, error)
	Update(id uint, req *dto.MenuUpdateReq) (*dto.MenuRes, error)
	Delete(id uint) error
}

type menuService struct {
	repo    repository.MenuRepository
	permSvc adminauth_service.PermissionService
}

func NewMenuService(repo repository.MenuRepository, permSvc adminauth_service.PermissionService) MenuService {
	return &menuService{repo: repo, permSvc: permSvc}
}

func mapToDto(m *entity.Menu) dto.MenuRes {
	return dto.MenuRes{
		ID:            m.ID,
		ParentID:      m.ParentID,
		Title:         m.Title,
		MenuKey:       m.MenuKey,
		Icon:          m.Icon,
		Url:           m.Url,
		PermissionBit: m.PermissionBit,
		IsActive:      m.IsActive,
		SortOrder:     m.SortOrder,
		CreatedAt:     m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     m.UpdatedAt.Format("2006-01-02 15:04:05"),
		Children:      []dto.MenuRes{},
	}
}

func buildTree(menus []entity.Menu, parentID uint) []dto.MenuRes {
	var tree []dto.MenuRes
	for _, m := range menus {
		if m.ParentID == parentID {
			node := mapToDto(&m)
			node.Children = buildTree(menus, m.ID)
			tree = append(tree, node)
		}
	}
	return tree
}

func (s *menuService) GetList() ([]dto.MenuRes, error) {
	menus, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var res []dto.MenuRes
	for _, m := range menus {
		res = append(res, mapToDto(&m))
	}
	if res == nil {
		res = []dto.MenuRes{}
	}
	return res, nil
}

func (s *menuService) GetTree() ([]dto.MenuRes, error) {
	menus, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	res := buildTree(menus, 0)
	if res == nil {
		res = []dto.MenuRes{}
	}
	return res, nil
}

func (s *menuService) GetMyMenuTree(adminID uint, isSuperAdmin bool) ([]dto.MenuRes, error) {
	allowedMenus, err := s.permSvc.GetAllowedMenus(adminID, isSuperAdmin)
	if err != nil {
		return nil, errors.New("không thể lấy quyền menu")
	}

	var activeMenus []entity.Menu
	for _, m := range allowedMenus {
		if m.IsActive {
			activeMenus = append(activeMenus, m)
		}
	}

	res := buildTree(activeMenus, 0)
	if res == nil {
		res = []dto.MenuRes{}
	}
	return res, nil
}

func (s *menuService) Create(req *dto.MenuCreateReq) (*dto.MenuRes, error) {
	if s.repo.CheckMenuKeyExists(req.MenuKey, 0) {
		return nil, errors.New("menu_key đã tồn tại")
	}

	if req.ParentID != 0 {
		if _, err := s.repo.FindByID(req.ParentID); err != nil {
			return nil, errors.New("parent_id không hợp lệ")
		}
	}

	m := &entity.Menu{
		ParentID:      req.ParentID,
		Title:         req.Title,
		MenuKey:       req.MenuKey,
		Icon:          req.Icon,
		Url:           req.Url,
		PermissionBit: req.PermissionBit,
		IsActive:      req.IsActive,
		SortOrder:     req.SortOrder,
	}

	if err := s.repo.Create(m); err != nil {
		return nil, errors.New("không thể tạo menu")
	}

	res := mapToDto(m)
	return &res, nil
}

func (s *menuService) Update(id uint, req *dto.MenuUpdateReq) (*dto.MenuRes, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy menu")
	}

	if s.repo.CheckMenuKeyExists(req.MenuKey, id) {
		return nil, errors.New("menu_key đã tồn tại")
	}

	if req.ParentID != 0 && req.ParentID != m.ParentID {
		if _, err := s.repo.FindByID(req.ParentID); err != nil {
			return nil, errors.New("parent_id không hợp lệ")
		}
		if req.ParentID == id {
			return nil, errors.New("parent_id không thể là chính nó")
		}
	}

	m.ParentID = req.ParentID
	m.Title = req.Title
	m.MenuKey = req.MenuKey
	m.Icon = req.Icon
	m.Url = req.Url
	m.PermissionBit = req.PermissionBit
	m.IsActive = req.IsActive
	m.SortOrder = req.SortOrder

	if err := s.repo.Update(m); err != nil {
		return nil, errors.New("lỗi cập nhật menu")
	}

	res := mapToDto(m)
	return &res, nil
}

func (s *menuService) Delete(id uint) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy menu")
	}

	childCount, err := s.repo.CountChildren(id)
	if err != nil {
		return errors.New("lỗi kiểm tra menu con")
	}
	if childCount > 0 {
		return errors.New("không thể xóa vì menu này đang có menu con")
	}

	if err := s.repo.Delete(m); err != nil {
		return errors.New("lỗi xóa menu")
	}
	return nil
}
