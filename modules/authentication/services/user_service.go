package services

import (
	"context"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
	ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error)
	CountUsers(ctx context.Context) (int64, error)
	SoftDeleteUser(ctx context.Context, user *models.User) error
	RestoreUser(ctx context.Context, user *models.User) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, user *models.User) (*models.User, error) {
	return s.userRepo.GetByID(ctx, user)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Delete(ctx, user)
}

func (s *userService) ListUsers(ctx context.Context, offset int, limit int) ([]*models.User, error) {
	return s.userRepo.List(ctx, offset, limit)
}

func (s *userService) CountUsers(ctx context.Context) (int64, error) {
	return s.userRepo.Count(ctx)
}

func (s *userService) SoftDeleteUser(ctx context.Context, user *models.User) error {
	return s.userRepo.SoftDelete(ctx, user)
}

func (s *userService) RestoreUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Restore(ctx, user)
}
