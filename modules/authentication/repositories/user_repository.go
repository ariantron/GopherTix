package repositories

import (
	"context"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authorization/constants"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, user *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
	List(ctx context.Context, offset, limit int, search *string) ([]*models.User, error)
	Count(ctx context.Context, search *string) (int64, error)
	SoftDelete(ctx context.Context, user *models.User) error
	Restore(ctx context.Context, user *models.User) error
}

type Repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetByID(ctx context.Context, user *models.User) (*models.User, error) {
	var result models.User
	err := r.db.WithContext(ctx).Unscoped().First(&result, user.ID).Error
	if err != nil {
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
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (r *Repository) Delete(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}

func (r *Repository) List(ctx context.Context, offset int, limit int, search *string) ([]*models.User, error) {
	query := r.db.WithContext(ctx)

	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
	}

	var users []*models.User
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) Count(ctx context.Context, search *string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.User{})
	if search != nil && *search != "" {
		query = query.Where("users.name ILIKE ?",
			constants.OwnerRoleID, "%"+*search+"%")
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *Repository) SoftDelete(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}

func (r *Repository) Restore(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Unscoped().Model(user).Update("deleted_at", nil).Error
}
