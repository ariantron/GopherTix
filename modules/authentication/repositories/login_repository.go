package repositories

import (
	"context"
	"gopher_tix/modules/authentication/models"
	errs "gopher_tix/packages/common/errors"
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
	if err := r.db.WithContext(ctx).Create(login).Error; err != nil {
		return errs.NewInternalServerError("Failed to create login record")
	}
	return nil
}
