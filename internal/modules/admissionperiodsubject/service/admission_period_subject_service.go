package service

import (
	"errors"

	"go_be_enrollment/internal/modules/admissionperiodsubject/dto"
	"go_be_enrollment/internal/modules/admissionperiodsubject/entity"
	"go_be_enrollment/internal/modules/admissionperiodsubject/repository"
	period_repo "go_be_enrollment/internal/modules/admissionperiod/repository"
	subject_repo "go_be_enrollment/internal/modules/subject/repository"
)

type AdmissionPeriodSubjectService interface {
	GetByAdmissionPeriodID(periodID uint) ([]dto.AdmissionPeriodSubjectRes, error)
	ReplaceSubjects(periodID uint, req *dto.ReplaceAdmissionPeriodSubjectsReq) ([]dto.AdmissionPeriodSubjectRes, error)
}

type admissionPeriodSubjectService struct {
	repo         repository.AdmissionPeriodSubjectRepository
	periodRepo   period_repo.AdmissionPeriodRepository
	subjectRepo  subject_repo.SubjectRepository
}

func NewAdmissionPeriodSubjectService(repo repository.AdmissionPeriodSubjectRepository, periodRepo period_repo.AdmissionPeriodRepository, subjectRepo subject_repo.SubjectRepository) AdmissionPeriodSubjectService {
	return &admissionPeriodSubjectService{
		repo:         repo,
		periodRepo:   periodRepo,
		subjectRepo:  subjectRepo,
	}
}

func (s *admissionPeriodSubjectService) GetByAdmissionPeriodID(periodID uint) ([]dto.AdmissionPeriodSubjectRes, error) {
	_, err := s.periodRepo.FindByID(periodID)
	if err != nil {
		return nil, errors.New("kỳ tuyển sinh không tồn tại")
	}

	subjects, err := s.repo.GetByAdmissionPeriodID(periodID)
	if err != nil {
		return nil, err
	}

	var res []dto.AdmissionPeriodSubjectRes
	for _, sub := range subjects {
		code := ""
		name := ""
		if sub.Subject != nil {
			code = sub.Subject.Code
			name = sub.Subject.Name
		}
		res = append(res, dto.AdmissionPeriodSubjectRes{
			ID:                sub.ID,
			AdmissionPeriodID: sub.AdmissionPeriodID,
			SubjectID:         sub.SubjectID,
			SubjectCode:       code,
			SubjectName:       name,
			Weight:            sub.Weight,
			IsRequired:        sub.IsRequired,
		})
	}
	if res == nil {
		res = []dto.AdmissionPeriodSubjectRes{}
	}
	return res, nil
}

func (s *admissionPeriodSubjectService) ReplaceSubjects(periodID uint, req *dto.ReplaceAdmissionPeriodSubjectsReq) ([]dto.AdmissionPeriodSubjectRes, error) {
	_, err := s.periodRepo.FindByID(periodID)
	if err != nil {
		return nil, errors.New("kỳ tuyển sinh không tồn tại")
	}

	subjectMap := make(map[uint]bool)
	var entities []entity.AdmissionPeriodSubject

	for _, item := range req.Subjects {
		if subjectMap[item.SubjectID] {
			return nil, errors.New("danh sách môn học có trùng lặp")
		}
		subjectMap[item.SubjectID] = true
		
		if item.Weight <= 0 {
			return nil, errors.New("trọng số phải lớn hơn 0")
		}

		_, err := s.subjectRepo.FindByID(item.SubjectID)
		if err != nil {
			return nil, errors.New("môn học không tồn tại")
		}

		entities = append(entities, entity.AdmissionPeriodSubject{
			AdmissionPeriodID: periodID,
			SubjectID:         item.SubjectID,
			Weight:            item.Weight,
			IsRequired:        item.IsRequired,
		})
	}

	if err := s.repo.ReplaceSubjects(periodID, entities); err != nil {
		return nil, err
	}

	return s.GetByAdmissionPeriodID(periodID)
}
