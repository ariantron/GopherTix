package repositories

import (
	"context"
	"gopher_tix/modules/authentication/models"
	"gorm.io/gorm"
)

type LoginRepository interface {
	CreateLoginRecord(ctx context.Context, login *models.Login) error
}

type loginRepository struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) LoginRepository {
	return &loginRepository{
		db: db,
	}
}

func (r *loginRepository) CreateLoginRecord(ctx context.Context, login *models.Login) error {
	return r.db.WithContext(ctx).Create(login).Error
}
