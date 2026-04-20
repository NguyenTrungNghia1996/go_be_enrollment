package service

import (
	"errors"

	"go_be_enrollment/internal/modules/academicrecord/dto"
	"go_be_enrollment/internal/modules/academicrecord/entity"
	"go_be_enrollment/internal/modules/academicrecord/repository"
)

type AcademicRecordService interface {
	GetAdminList(appID uint) ([]dto.AcademicRecordRes, error)
	GetUserList(appID, userID uint) ([]dto.AcademicRecordRes, error)
	Create(appID, userID uint, req *dto.AcademicRecordReq) (*dto.AcademicRecordRes, error)
	Update(id, userID uint, req *dto.AcademicRecordReq) (*dto.AcademicRecordRes, error)
	Delete(id, userID uint) error
}

type academicRecordService struct {
	repo repository.AcademicRecordRepository
}

func NewAcademicRecordService(repo repository.AcademicRecordRepository) AcademicRecordService {
	return &academicRecordService{repo: repo}
}

func mapToDto(rec *entity.AcademicRecord) dto.AcademicRecordRes {
	return dto.AcademicRecordRes{
		ID:                  rec.ID,
		ApplicationID:       rec.ApplicationID,
		GradeLevel:          rec.GradeLevel,
		SchoolName:          rec.SchoolName,
		SchoolLocation:      rec.SchoolLocation,
		SchoolYear:          rec.SchoolYear,
		AcademicPerformance: rec.AcademicPerformance,
		ConductRating:       rec.ConductRating,
		CreatedAt:           rec.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           rec.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *academicRecordService) GetAdminList(appID uint) ([]dto.AcademicRecordRes, error) {
	list, err := s.repo.GetByApplicationID(appID)
	if err != nil {
		return nil, errors.New("lỗi lấy học bạ")
	}

	res := make([]dto.AcademicRecordRes, 0)
	for _, rec := range list {
		res = append(res, mapToDto(&rec))
	}
	return res, nil
}

func (s *academicRecordService) GetUserList(appID, userID uint) ([]dto.AcademicRecordRes, error) {
	isOwner, _, err := s.repo.VerifyApplicationAccess(appID, userID)
	if err != nil {
		return nil, errors.New("lỗi hệ thống khi xác thực hồ sơ")
	}
	if !isOwner {
		// Dù hồ sơ tồn tại nhưng user không sở hữu, trả chung lỗi không tìm thấy
		return nil, errors.New("không tìm thấy hồ sơ của bạn")
	}

	list, err := s.repo.GetByApplicationID(appID)
	if err != nil {
		return nil, errors.New("lỗi lấy học bạ")
	}

	res := make([]dto.AcademicRecordRes, 0)
	for _, rec := range list {
		res = append(res, mapToDto(&rec))
	}
	return res, nil
}

func (s *academicRecordService) Create(appID, userID uint, req *dto.AcademicRecordReq) (*dto.AcademicRecordRes, error) {
	isOwner, status, err := s.repo.VerifyApplicationAccess(appID, userID)
	if err != nil {
		return nil, errors.New("lỗi hệ thống khi xác thực hồ sơ")
	}
	if !isOwner {
		return nil, errors.New("không tìm thấy hồ sơ của bạn")
	}
	if status != "Draft" {
		return nil, errors.New("chỉ có thể thêm học bạ khi hồ sơ ở trạng thái nháp")
	}

	if s.repo.CheckDuplicateGrade(appID, req.GradeLevel, 0) {
		return nil, errors.New("học bạ cho khối lớp này đã tồn tại")
	}

	record := &entity.AcademicRecord{
		ApplicationID:       appID,
		GradeLevel:          req.GradeLevel,
		SchoolName:          req.SchoolName,
		SchoolLocation:      req.SchoolLocation,
		SchoolYear:          req.SchoolYear,
		AcademicPerformance: req.AcademicPerformance,
		ConductRating:       req.ConductRating,
	}

	if err := s.repo.Create(record); err != nil {
		return nil, errors.New("lỗi khi lưu học bạ")
	}

	res := mapToDto(record)
	return &res, nil
}

func (s *academicRecordService) Update(id, userID uint, req *dto.AcademicRecordReq) (*dto.AcademicRecordRes, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("không tìm thấy học bạ")
	}

	isOwner, status, err := s.repo.VerifyApplicationAccess(record.ApplicationID, userID)
	if err != nil {
		return nil, errors.New("lỗi hệ thống khi xác thực hồ sơ")
	}
	if !isOwner {
		return nil, errors.New("bạn không có quyền chỉnh sửa học bạ này")
	}
	if status != "Draft" {
		return nil, errors.New("chỉ có thể sửa học bạ khi hồ sơ ở trạng thái nháp")
	}

	if s.repo.CheckDuplicateGrade(record.ApplicationID, req.GradeLevel, record.ID) {
		return nil, errors.New("học bạ cho khối/lớp này đã tồn tại")
	}

	record.GradeLevel = req.GradeLevel
	record.SchoolName = req.SchoolName
	record.SchoolLocation = req.SchoolLocation
	record.SchoolYear = req.SchoolYear
	record.AcademicPerformance = req.AcademicPerformance
	record.ConductRating = req.ConductRating

	if err := s.repo.Update(record); err != nil {
		return nil, errors.New("lỗi cập nhật học bạ")
	}

	res := mapToDto(record)
	return &res, nil
}

func (s *academicRecordService) Delete(id, userID uint) error {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("không tìm thấy học bạ")
	}

	isOwner, status, err := s.repo.VerifyApplicationAccess(record.ApplicationID, userID)
	if err != nil {
		return errors.New("lỗi hệ thống khi xác thực hồ sơ")
	}
	if !isOwner {
		return errors.New("bạn không có quyền xóa học bạ này")
	}
	if status != "Draft" {
		return errors.New("chỉ có thể xóa học bạ khi hồ sơ ở trạng thái nháp")
	}

	if err := s.repo.Delete(id); err != nil {
		return errors.New("lỗi khi xóa học bạ")
	}

	return nil
}
