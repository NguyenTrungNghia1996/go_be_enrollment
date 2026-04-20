package dto

type AcademicRecordReq struct {
	GradeLevel          int    `json:"grade_level" validate:"required,min=1,max=5"`
	SchoolName          string `json:"school_name" validate:"required"`
	SchoolLocation      string `json:"school_location"`
	SchoolYear          string `json:"school_year"`
	AcademicPerformance string `json:"academic_performance"`
	ConductRating       string `json:"conduct_rating"`
}

type AcademicRecordRes struct {
	ID                  uint   `json:"id"`
	ApplicationID       uint   `json:"application_id"`
	GradeLevel          int    `json:"grade_level"`
	SchoolName          string `json:"school_name"`
	SchoolLocation      string `json:"school_location"`
	SchoolYear          string `json:"school_year"`
	AcademicPerformance string `json:"academic_performance"`
	ConductRating       string `json:"conduct_rating"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
