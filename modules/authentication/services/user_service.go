package services

import (
	"context"
	"github.com/google/uuid"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/repositories"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
	List(ctx context.Context, offset, limit int, search *string) ([]*models.User, error)
	Count(ctx context.Context, search *string) (int64, error)
	Deactivate(ctx context.Context, user *models.User) error
	Activate(ctx context.Context, user *models.User) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, user *models.User) error {
	return s.userRepo.Delete(ctx, user)
}

func (s *userService) List(ctx context.Context, offset int, limit int, search *string) ([]*models.User, error) {
	return s.userRepo.List(ctx, offset, limit, search)
}

func (s *userService) Count(ctx context.Context, search *string) (int64, error) {
	return s.userRepo.Count(ctx, search)
}

func (s *userService) Deactivate(ctx context.Context, user *models.User) error {
	return s.userRepo.SoftDelete(ctx, user)
}

func (s *userService) Activate(ctx context.Context, user *models.User) error {
	return s.userRepo.Restore(ctx, user)
}
