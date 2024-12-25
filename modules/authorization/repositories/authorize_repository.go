package repositories

import (
	"context"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/constants"
	"gopher_tix/modules/authorization/models"
	"gorm.io/gorm"
)

type AuthorizeRepository struct {
	db *gorm.DB
}

func NewAuthorizeRepository(db *gorm.DB) *AuthorizeRepository {
	return &AuthorizeRepository{db: db}
}

func (r *AuthorizeRepository) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).Model(&models.UserRole{}).
		Select("1").
		Where("user_id = ? AND role_id = ? AND group_id IS NULL", userID, constants.AdminRoleID).
		Limit(1).
		Find(&exists).Error
	return exists, err
}

func (r *AuthorizeRepository) HasRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID int8) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).Model(&models.UserRole{}).
		Select("1").
		Where("user_id = ? AND group_id = ? AND role_id = ?", userID, groupID, roleID).
		Limit(1).
		Find(&exists).Error
	return exists, err
}

func (r *AuthorizeRepository) GetUserRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.WithContext(ctx).Preload("Role").
		Where("user_id = ? AND group_id = ?", userID, groupID).
		First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r *AuthorizeRepository) AssignRole(ctx context.Context, userRole *models.UserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *AuthorizeRepository) UnassignRole(ctx context.Context, userRole *models.UserRole) error {
	return r.db.WithContext(ctx).Delete(userRole).Error
}
