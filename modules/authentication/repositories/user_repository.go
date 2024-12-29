package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gopher_tix/modules/authentication/models"
	errs "gopher_tix/packages/common/errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
	List(ctx context.Context, offset, limit int, search *string) ([]*models.User, error)
	Count(ctx context.Context, search *string) (int64, error)
	SoftDelete(ctx context.Context, user *models.User) error
	Restore(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Unscoped().First(&user, id).Error; err != nil {
		return nil, handleError(err, "User ID "+id.String())
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, handleError(err, "User")
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return handleError(err, "Failed to update user")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Delete(user).Error; err != nil {
		return handleError(err, "Failed to delete user")
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int, search *string) ([]*models.User, error) {
	query := r.db.WithContext(ctx)
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
	}

	var users []*models.User
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, errs.NewInternalServerError("Failed to list users")
	}
	return users, nil
}

func (r *userRepository) Count(ctx context.Context, search *string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.User{})
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, errs.NewInternalServerError("Failed to count users")
	}
	return count, nil
}

func (r *userRepository) SoftDelete(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Delete(user).Error; err != nil {
		return handleError(err, "Failed to deactivate user")
	}
	return nil
}

func (r *userRepository) Restore(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Unscoped().Model(user).Update("deleted_at", nil).Error; err != nil {
		return handleError(err, "Failed to activate user")
	}
	return nil
}

func handleError(err error, message string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewNotFoundError("User")
	}
	return errs.NewInternalServerError(message)
}
