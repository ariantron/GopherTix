package repositories

import (
	"context"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/constants"
	"gopher_tix/modules/authorization/models"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) List(ctx context.Context, offset, limit int) ([]models.Group, error) {
	var groups []models.Group
	err := r.db.WithContext(ctx).
		Preload("ParentGroup").
		Preload("SubGroups").
		Offset(offset).
		Limit(limit).
		Find(&groups).Error
	return groups, err
}

func (r *GroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.WithContext(ctx).
		Preload("ParentGroup").
		Preload("SubGroups").
		First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) Create(ctx context.Context, group *models.Group, ownerUserID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		userRole := models.UserRole{
			UserID:  ownerUserID,
			RoleID:  constants.OwnerRoleID,
			GroupID: group.ID,
		}
		return tx.Create(&userRole).Error
	})
}

func (r *GroupRepository) Update(ctx context.Context, group *models.Group) error {
	return r.db.WithContext(ctx).Model(group).Updates(map[string]interface{}{
		"name": group.Name,
	}).Error
}

func (r *GroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Group{}, id).Error
}

func (r *GroupRepository) GetOwner(ctx context.Context, groupID uuid.UUID) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.WithContext(ctx).Where("group_id = ? AND role_id = ?", groupID, constants.OwnerRoleID).
		First(&userRole).Error
	return &userRole, err
}

func (r *GroupRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&userRoles).Error
	return userRoles, err
}
