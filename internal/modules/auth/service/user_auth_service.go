package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"go_be_enrollment/internal/config"
	"go_be_enrollment/internal/modules/auth/dto"
	"go_be_enrollment/internal/modules/auth/entity"
	"go_be_enrollment/internal/modules/auth/repository"
	"go_be_enrollment/pkg/logger"
	"go_be_enrollment/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
	Activate(ctx context.Context, req *dto.ActivateRequest) error
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

	var token *string
	var otp *string
	var expiresAt *time.Time
	isActive := true

	if req.Email != "" {
		isActive = false
		t := uuid.New().String()
		token = &t

		n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
		o := fmt.Sprintf("%06d", n.Int64())
		otp = &o

		expr := time.Now().Add(24 * time.Hour)
		expiresAt = &expr
	}

	user := &entity.UserAccount{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		IsActive:     isActive,
		ActivationToken:     token,
		ActivationOTP:       otp,
		ActivationExpiresAt: expiresAt,
	}

	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = &req.PhoneNumber
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	if req.Email != "" {
		go func(email, fullName, t, o string) {
			smtpCfg := utils.SMTPConfig{
				Host:     s.cfg.SMTPHost,
				Port:     s.cfg.SMTPPort,
				User:     s.cfg.SMTPUser,
				Password: s.cfg.SMTPPassword,
				From:     s.cfg.SMTPFrom,
			}

			if smtpCfg.Host == "" {
				return
			}

			subject := "Xác nhận đăng ký tài khoản và mã Kích hoạt"
			body := fmt.Sprintf("Chào %s,\n\nBạn đã đăng ký tài khoản thành công.\nMã OTP kích hoạt của bạn là: %s\nHoặc có thể gửi mã Token này: %s lên hệ thống để kích hoạt!\n\nTrân trọng,\nĐội ngũ quản trị", fullName, o, t)

			err := utils.SendEmail(smtpCfg, []string{email}, subject, body)
			if err != nil {
				logger.Log.Error("Failed to send confirmation email", zap.Error(err), zap.String("email", email))
			}
		}(req.Email, req.FullName, *user.ActivationToken, *user.ActivationOTP)
	}

	return nil
}

func (s *userAuthService) Activate(ctx context.Context, req *dto.ActivateRequest) error {
	if req.OTP == "" && req.Token == "" {
		return errors.New("vui lòng cung cấp otp hoặc token")
	}

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return ErrUserNotFound
	}

	if user.IsActive {
		return nil
	}

	if user.ActivationExpiresAt != nil && time.Now().After(*user.ActivationExpiresAt) {
		return errors.New("mã kích hoạt đã hết hạn")
	}

	if req.OTP != "" && (user.ActivationOTP == nil || *user.ActivationOTP != req.OTP) {
		return errors.New("mã otp không chính xác")
	}

	if req.Token != "" && (user.ActivationToken == nil || *user.ActivationToken != req.Token) {
		return errors.New("token không chính xác")
	}

	user.IsActive = true
	user.ActivationOTP = nil
	user.ActivationToken = nil
	user.ActivationExpiresAt = nil

	return s.repo.Update(ctx, user)
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
