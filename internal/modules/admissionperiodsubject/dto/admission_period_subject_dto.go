package dto

type AdmissionPeriodSubjectItemReq struct {
	SubjectID  uint    `json:"subject_id" validate:"required"`
	Weight     float64 `json:"weight" validate:"required,gt=0"`
	IsRequired bool    `json:"is_required"`
}

type ReplaceAdmissionPeriodSubjectsReq struct {
	Subjects []AdmissionPeriodSubjectItemReq `json:"subjects"`
}

type AdmissionPeriodSubjectRes struct {
	ID                uint    `json:"id"`
	AdmissionPeriodID uint    `json:"admission_period_id"`
	SubjectID         uint    `json:"subject_id"`
	SubjectCode       string  `json:"subject_code"`
	SubjectName       string  `json:"subject_name"`
	Weight            float64 `json:"weight"`
	IsRequired        bool    `json:"is_required"`
}
