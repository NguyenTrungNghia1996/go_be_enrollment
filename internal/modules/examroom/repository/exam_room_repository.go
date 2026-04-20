package repository

import (
	"go_be_enrollment/internal/modules/examroom/dto"
	"go_be_enrollment/internal/modules/examroom/entity"

	"gorm.io/gorm"
)

type ExamRoomRepository interface {
	GetList(filter *dto.ExamRoomFilter) ([]entity.ExamRoom, int64, error)
	FindByID(id uint) (*entity.ExamRoom, error)
	Create(e *entity.ExamRoom) error
	Update(e *entity.ExamRoom) error
	Delete(id uint) error
	HasRelatedAssignments(id uint) (bool, error)
}

type examRoomRepository struct {
	db *gorm.DB
}

func NewExamRoomRepository(db *gorm.DB) ExamRoomRepository {
	return &examRoomRepository{db: db}
}

func (r *examRoomRepository) GetList(filter *dto.ExamRoomFilter) ([]entity.ExamRoom, int64, error) {
	query := r.db.Model(&entity.ExamRoom{})

	if filter.Keyword != "" {
		query = query.Where("room_name LIKE ? OR location LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}
	offset := (page - 1) * limit

	var rooms []entity.ExamRoom
	if err := query.Order("id desc").Offset(offset).Limit(limit).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *examRoomRepository) FindByID(id uint) (*entity.ExamRoom, error) {
	var room entity.ExamRoom
	if err := r.db.First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *examRoomRepository) Create(e *entity.ExamRoom) error {
	return r.db.Create(e).Error
}

func (r *examRoomRepository) Update(e *entity.ExamRoom) error {
	return r.db.Save(e).Error
}

func (r *examRoomRepository) Delete(id uint) error {
	return r.db.Delete(&entity.ExamRoom{}, id).Error
}

func (r *examRoomRepository) HasRelatedAssignments(id uint) (bool, error) {
	if r.db.Migrator().HasTable("exam_room_assignments") {
		var count int64
		if err := r.db.Table("exam_room_assignments").Where("exam_room_id = ?", id).Count(&count).Error; err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
	}

	if r.db.Migrator().HasTable("examiner_assignments") {
		var count int64
		if err := r.db.Table("examiner_assignments").Where("exam_room_id = ?", id).Count(&count).Error; err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}
