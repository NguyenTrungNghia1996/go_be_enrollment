package dto

type ExamRoomFilter struct {
	Keyword  string `query:"keyword"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

type ExamRoomCreateReq struct {
	RoomName string `json:"room_name" validate:"required"`
	Location string `json:"location" validate:"required"`
	Capacity int    `json:"capacity"`
}

type ExamRoomUpdateReq struct {
	RoomName string `json:"room_name" validate:"required"`
	Location string `json:"location" validate:"required"`
	Capacity int    `json:"capacity"`
}

type ExamRoomRes struct {
	ID        uint   `json:"id"`
	RoomName  string `json:"room_name"`
	Location  string `json:"location"`
	Capacity  int    `json:"capacity"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PaginatedExamRoomRes struct {
	Data       []ExamRoomRes `json:"data"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}
