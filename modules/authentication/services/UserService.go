package services

import (
	"context"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/repositories"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
	ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error)
	CountUsers(ctx context.Context) (int64, error)
	UpdateUserVerificationStatus(ctx context.Context, user *models.User, verifiedAt time.Time) error
	SoftDeleteUser(ctx context.Context, user *models.User) error
	RestoreUser(ctx context.Context, user *models.User) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, user *models.User) (*models.User, error) {
	return s.repo.GetByID(ctx, user)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, user *models.User) error {
	return s.repo.Delete(ctx, user)
}

func (s *userService) ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *userService) CountUsers(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *userService) UpdateUserVerificationStatus(ctx context.Context, user *models.User, verifiedAt time.Time) error {
	return s.repo.UpdateVerificationStatus(ctx, user, verifiedAt)
}

func (s *userService) SoftDeleteUser(ctx context.Context, user *models.User) error {
	return s.repo.SoftDelete(ctx, user)
}

func (s *userService) RestoreUser(ctx context.Context, user *models.User) error {
	return s.repo.Restore(ctx, user)
}
