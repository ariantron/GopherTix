package repositories

import (
	"context"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/constants"
	"gopher_tix/modules/authorization/models"
	"gorm.io/gorm"
)

type AuthorizeRepository interface {
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	HasRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) (bool, error)
	GetUserRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (*models.UserRole, error)
	AssignRole(ctx context.Context, userRole *models.UserRole) error
	UnassignRole(ctx context.Context, userRole *models.UserRole) error
}

type authorizeRepository struct {
	db *gorm.DB
}

func NewAuthorizeRepository(db *gorm.DB) AuthorizeRepository {
	return &authorizeRepository{db: db}
}

func (r *authorizeRepository) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).Model(&models.UserRole{}).
		Select("1").
		Where("user_id = ? AND role_id = ? AND group_id IS NULL", userID, constants.AdminRoleID).
		Limit(1).
		Find(&exists).Error
	return exists, err
}

func (r *authorizeRepository) HasRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).Model(&models.UserRole{}).
		Select("1").
		Where("user_id = ? AND group_id = ? AND role_id = ?", userID, groupID, roleID).
		Limit(1).
		Find(&exists).Error
	return exists, err
}

func (r *authorizeRepository) GetUserRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.WithContext(ctx).Preload("Role").
		Where("user_id = ? AND group_id = ?", userID, groupID).
		First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r *authorizeRepository) AssignRole(ctx context.Context, userRole *models.UserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *authorizeRepository) UnassignRole(ctx context.Context, userRole *models.UserRole) error {
	return r.db.WithContext(ctx).Delete(userRole).Error
}
