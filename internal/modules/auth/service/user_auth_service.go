package service

import (
	"context"
	"errors"

	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/modules/auth/dto"
	"go_be_enrollment/internal/modules/auth/entity"
	"go_be_enrollment/internal/modules/auth/repository"
	"go_be_enrollment/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists    = errors.New("email already exists")
	ErrUserInactive   = errors.New("user account is inactive")
	ErrInvalidLogin   = errors.New("invalid username or password")
	ErrUserNotFound   = errors.New("user not found")
)

type UserAuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error)
	GetMe(ctx context.Context, userID uint) (*dto.UserInfoResponse, error)
}

type userAuthService struct {
	repo repository.UserAccountRepository
	cfg  *config.Config
}

func NewUserAuthService(repo repository.UserAccountRepository, cfg *config.Config) UserAuthService {
	return &userAuthService{repo: repo, cfg: cfg}
}

func (s *userAuthService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	// Check username
	_, err := s.repo.FindByUsername(ctx, req.Username)
	if err == nil {
		return ErrUsernameExists
	}

	// Check email optionally
	if req.Email != "" {
		_, err = s.repo.FindByEmail(ctx, req.Email)
		if err == nil {
			return ErrEmailExists
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.UserAccount{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		IsActive:     true,
	}

	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = &req.PhoneNumber
	}

	return s.repo.Create(ctx, user)
}

func (s *userAuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, ErrInvalidLogin
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidLogin
	}

	token, err := utils.GenerateUserToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTExpiresIn)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{AccessToken: token}, nil
}

func (s *userAuthService) GetMe(ctx context.Context, userID uint) (*dto.UserInfoResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &dto.UserInfoResponse{
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
	}, nil
}
