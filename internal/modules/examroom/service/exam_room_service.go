package service

import (
	"errors"
	"math"
	"time"

	"go_be_enrollment/internal/modules/examroom/dto"
	"go_be_enrollment/internal/modules/examroom/entity"
	"go_be_enrollment/internal/modules/examroom/repository"
)

type ExamRoomService interface {
	GetList(filter *dto.ExamRoomFilter) (*dto.PaginatedExamRoomRes, error)
	GetDetail(id uint) (*dto.ExamRoomRes, error)
	Create(req *dto.ExamRoomCreateReq) (*dto.ExamRoomRes, error)
	Update(id uint, req *dto.ExamRoomUpdateReq) (*dto.ExamRoomRes, error)
	Delete(id uint) error
}

type examRoomService struct {
	repo repository.ExamRoomRepository
}

func NewExamRoomService(repo repository.ExamRoomRepository) ExamRoomService {
	return &examRoomService{repo: repo}
}

func (s *examRoomService) GetList(filter *dto.ExamRoomFilter) (*dto.PaginatedExamRoomRes, error) {
	rooms, total, err := s.repo.GetList(filter)
	if err != nil {
		return nil, err
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	var resData []dto.ExamRoomRes
	for _, r := range rooms {
		resData = append(resData, dto.ExamRoomRes{
			ID:        r.ID,
			RoomName:  r.RoomName,
			Location:  r.Location,
			Capacity:  r.Capacity,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
			UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
		})
	}
	
	if resData == nil {
		resData = []dto.ExamRoomRes{}
	}

	return &dto.PaginatedExamRoomRes{
		Data:       resData,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *examRoomService) GetDetail(id uint) (*dto.ExamRoomRes, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("phòng thi không tồn tại")
	}

	return &dto.ExamRoomRes{
		ID:        room.ID,
		RoomName:  room.RoomName,
		Location:  room.Location,
		Capacity:  room.Capacity,
		CreatedAt: room.CreatedAt.Format(time.RFC3339),
		UpdatedAt: room.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *examRoomService) Create(req *dto.ExamRoomCreateReq) (*dto.ExamRoomRes, error) {
	if req.Capacity < 0 {
		return nil, errors.New("sức chứa không được âm")
	}

	room := &entity.ExamRoom{
		RoomName: req.RoomName,
		Location: req.Location,
		Capacity: req.Capacity,
	}

	if err := s.repo.Create(room); err != nil {
		return nil, err
	}

	return &dto.ExamRoomRes{
		ID:        room.ID,
		RoomName:  room.RoomName,
		Location:  room.Location,
		Capacity:  room.Capacity,
		CreatedAt: room.CreatedAt.Format(time.RFC3339),
		UpdatedAt: room.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *examRoomService) Update(id uint, req *dto.ExamRoomUpdateReq) (*dto.ExamRoomRes, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("phòng thi không tồn tại")
	}

	if req.Capacity < 0 {
		return nil, errors.New("sức chứa không được âm")
	}

	room.RoomName = req.RoomName
	room.Location = req.Location
	room.Capacity = req.Capacity

	if err := s.repo.Update(room); err != nil {
		return nil, err
	}

	return &dto.ExamRoomRes{
		ID:        room.ID,
		RoomName:  room.RoomName,
		Location:  room.Location,
		Capacity:  room.Capacity,
		CreatedAt: room.CreatedAt.Format(time.RFC3339),
		UpdatedAt: room.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *examRoomService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("phòng thi không tồn tại")
	}

	hasRelated, err := s.repo.HasRelatedAssignments(id)
	if err != nil {
		return err
	}
	if hasRelated {
		return errors.New("không thể xóa phòng thi vì đã có dữ liệu phân công liên quan")
	}

	return s.repo.Delete(id)
}
