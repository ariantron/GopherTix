package repositories

import (
	"context"
	"gopher_tix/modules/authentication/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, user *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)
	SoftDelete(ctx context.Context, user *models.User) error
	Restore(ctx context.Context, user *models.User) error
}

type Repository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetByID(ctx context.Context, user *models.User) (*models.User, error) {
	var result models.User
	if err := r.db.WithContext(ctx).First(&result, "id = ?", user.ID).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var result models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *Repository) Delete(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) SoftDelete(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}

func (r *Repository) Restore(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Unscoped().Model(user).Update("deleted_at", nil).Error
}
