package interfaces

import (
	"context"
	"gopher_tix/modules/authentication/models"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, user *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error

	List(ctx context.Context, offset, limit int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)

	UpdateVerificationStatus(ctx context.Context, user *models.User, verifiedAt time.Time) error
	GetUnverifiedUsers(ctx context.Context) ([]*models.User, error)

	SoftDelete(ctx context.Context, user *models.User) error
	Restore(ctx context.Context, user *models.User) error
}
