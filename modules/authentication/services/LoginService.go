package services

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopher_tix/configs"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/repositories"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
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
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.generateJWTToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *loginService) generateJWTToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(configs.SecretKey)
}
