package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/constants"
	"gopher_tix/modules/authorization/models"
	errs "gopher_tix/packages/common/errors"
	"gorm.io/gorm"
)

type GroupRepository interface {
	List(ctx context.Context, offset, limit int, search *string) ([]models.Group, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error)
	Create(ctx context.Context, group *models.Group, ownerUserID uuid.UUID) error
	Update(ctx context.Context, group *models.Group, ownerUserID *uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetOwner(ctx context.Context, groupID uuid.UUID) (*models.UserRole, error)
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]models.UserRole, error)
	Count(ctx context.Context, search *string) (int64, error)
	AddMember(ctx context.Context, userRole *models.UserRole) error
	RemoveMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) error
}

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) List(ctx context.Context, offset, limit int, search *string) ([]models.Group, error) {
	var groups []models.Group
	query := r.db.WithContext(ctx)
	if search != nil && *search != "" {
		query = query.Joins("JOIN user_roles ON user_roles.group_id = groups.id").
			Joins("JOIN users ON users.id = user_roles.user_id").
			Where("user_roles.role_id = ? AND users.name ILIKE ?",
				constants.OwnerRoleID, "%"+*search+"%")
	}
	err := query.Preload("ParentGroup").
		Preload("SubGroups").
		Offset(offset).
		Limit(limit).
		Find(&groups).Error
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to list groups")
	}
	return groups, nil
}

func (r *groupRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.WithContext(ctx).
		Preload("ParentGroup").
		Preload("SubGroups").
		First(&group, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("Group")
		}
		return nil, errs.NewInternalServerError("Failed to get group")
	}
	return &group, nil
}

func (r *groupRepository) Create(ctx context.Context, group *models.Group, ownerUserID uuid.UUID) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
	if err != nil {
		return errs.NewInternalServerError("Failed to create group")
	}
	return nil
}

func (r *groupRepository) Update(ctx context.Context, group *models.Group, ownerUserID *uuid.UUID) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(group).Updates(map[string]interface{}{
			"name": group.Name,
		}).Error; err != nil {
			return err
		}
		if ownerUserID != nil {
			if err := tx.Where("group_id = ? AND role_id = ?", group.ID, constants.OwnerRoleID).
				Delete(&models.UserRole{}).Error; err != nil {
				return errs.NewInternalServerError("Failed to change group owner")
			}
			newOwnerRole := models.UserRole{
				UserID:  *ownerUserID,
				RoleID:  constants.OwnerRoleID,
				GroupID: group.ID,
			}
			if err := tx.Create(&newOwnerRole).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errs.NewInternalServerError("Failed to update group")
	}
	return nil
}

func (r *groupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&models.Group{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewNotFoundError("Group")
		}
		return errs.NewInternalServerError("Failed to delete group")
	}
	return nil
}

func (r *groupRepository) GetOwner(ctx context.Context, groupID uuid.UUID) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.WithContext(ctx).Where("group_id = ? AND role_id = ?", groupID, constants.OwnerRoleID).
		First(&userRole).Error
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to get group owner")
	}
	return &userRole, nil
}

func (r *groupRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&userRoles).Error
	if err != nil {
		return nil, errs.NewInternalServerError("Failed to get group members")
	}
	return userRoles, err
}

func (r *groupRepository) Count(ctx context.Context, search *string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Group{})

	if search != nil && *search != "" {
		query = query.Joins("JOIN user_roles ON user_roles.group_id = groups.id").
			Joins("JOIN users ON users.id = user_roles.user_id").
			Where("user_roles.role_id = ? AND users.name ILIKE ?",
				constants.OwnerRoleID, "%"+*search+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, errs.NewInternalServerError("Failed to count groups")
	}
	return count, nil
}

func (r *groupRepository) AddMember(ctx context.Context, userRole *models.UserRole) error {
	err := r.db.WithContext(ctx).Create(userRole).Error
	if err != nil {
		return errs.NewInternalServerError("Failed to add member to group")
	}
	return nil
}

func (r *groupRepository) RemoveMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&models.UserRole{}).Error
	if err != nil {
		return errs.NewInternalServerError("Failed to remove member from group")
	}
	return nil
}
