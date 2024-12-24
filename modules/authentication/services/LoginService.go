package services

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gopher_tix/modules/authentication/middlewares"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/repositories"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrFailedGenerateToken = errors.New("failed to generate token")
)

type LoginService interface {
	CreateLoginRecord(ctx context.Context, login *models.Login) error
	ValidateUserCredentials(ctx context.Context, email, password string) (*models.User, string, error)
}

type loginService struct {
	loginRepo   repositories.LoginRepository
	userService UserService
}

func NewLoginService(loginRepo repositories.LoginRepository, userService UserService) LoginService {
	return &loginService{
		loginRepo:   loginRepo,
		userService: userService,
	}
}

func (s *loginService) CreateLoginRecord(ctx context.Context, login *models.Login) error {
	return s.loginRepo.CreateLoginRecord(ctx, login)
}

func (s *loginService) ValidateUserCredentials(ctx context.Context, email, password string) (*models.User, string, error) {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, "", ErrInvalidCredentials
	}

	token, err := middlewares.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return user, "", ErrFailedGenerateToken
	}

	return user, token, nil
}
