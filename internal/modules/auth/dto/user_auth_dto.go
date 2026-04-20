package dto

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=4,max=50"`
	Password    string `json:"password" validate:"required,min=6"`
	FullName    string `json:"full_name" validate:"required"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserInfoResponse struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	FullName    string  `json:"full_name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	IsActive    bool    `json:"is_active"`
}

type ActivateRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token,omitempty"`
	OTP   string `json:"otp,omitempty"`
}
