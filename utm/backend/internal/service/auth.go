package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/utm/backend/internal/middleware"
	"github.com/utm/backend/internal/model"
	"github.com/utm/backend/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email        string
	Password     string
	FullName     string
	Phone        string
	Organization string
	LicenseNo    string
}

type LoginResponse struct {
	AccessToken string             `json:"access_token"`
	User        *model.UserProfile `json:"user"`
}

type AuthService struct {
	userRepo    *postgres.UserRepository
	jwtSecret   string
	expiryHours int
}

func NewAuthService(userRepo *postgres.UserRepository, jwtSecret string, expiryHours int) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		expiryHours: expiryHours,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*model.UserProfile, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: string(hash),
		FullName:     input.FullName,
		Phone:        input.Phone,
		Role:         model.RolePilot,
		Status:       model.UserStatusPending,
		Organization: input.Organization,
		LicenseNo:    input.LicenseNo,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &model.UserProfile{
		ID:           user.ID,
		Email:        user.Email,
		FullName:     user.FullName,
		Phone:        user.Phone,
		Role:         user.Role,
		Status:       user.Status,
		Organization: user.Organization,
		LicenseNo:    user.LicenseNo,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Status == model.UserStatusBanned {
		return nil, errors.New("account is banned")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	claims := &middleware.Claims{
		UserID: user.ID.String(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: signed,
		User: &model.UserProfile{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			Status:   user.Status,
		},
	}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.UserProfile{
		ID:           user.ID,
		Email:        user.Email,
		FullName:     user.FullName,
		Phone:        user.Phone,
		Role:         user.Role,
		Status:       user.Status,
		Organization: user.Organization,
		LicenseNo:    user.LicenseNo,
	}, nil
}
