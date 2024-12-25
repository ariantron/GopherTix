package repositories

import (
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/models"
	"gorm.io/gorm"
)

type AuthorizeRepository struct {
	db *gorm.DB
}

func NewAuthorizeRepository(db *gorm.DB) *AuthorizeRepository {
	return &AuthorizeRepository{db: db}
}

func (r *AuthorizeRepository) IsAdmin(userID uuid.UUID) (bool, error) {
	var userRole models.UserRole
	err := r.db.Where("user_id = ? AND group_id IS NULL", userID).First(&userRole).Error
	return err == nil, nil
}

func (r *AuthorizeRepository) GetUserRole(userID uuid.UUID, groupID uuid.UUID) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r *AuthorizeRepository) AssignRole(userRole *models.UserRole) error {
	return r.db.Create(userRole).Error
}

func (r *AuthorizeRepository) UnassignRole(userRole *models.UserRole) error {
	return r.db.Delete(userRole).Error
}

func (r *AuthorizeRepository) GetGroupOwner(groupID uuid.UUID) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.Joins("Role").
		Where("group_id = ? AND roles.name = ?", groupID, "Owner").
		First(&userRole).Error
	return &userRole, err
}

func (r *AuthorizeRepository) GetGroupMembers(groupID uuid.UUID) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	err := r.db.Where("group_id = ?", groupID).Find(&userRoles).Error
	return userRoles, err
}
