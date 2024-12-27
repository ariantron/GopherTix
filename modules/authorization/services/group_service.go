package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	autzerrors "gopher_tix/modules/authorization/errors"
	"gopher_tix/modules/authorization/models"
	"gopher_tix/modules/authorization/repositories"
	"gorm.io/gorm"
)

type GroupService interface {
	List(ctx context.Context, offset, limit int, search *string) ([]models.Group, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error)
	Create(ctx context.Context, group *models.Group, ownerUserID uuid.UUID, currentUserID uuid.UUID) error
	Update(ctx context.Context, group *models.Group, ownerUserID *uuid.UUID, currentUserID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error
	GetMembers(ctx context.Context, groupID uuid.UUID, currentUserID uuid.UUID) ([]models.UserRole, error)
	Count(ctx context.Context, search *string) (int64, error)
	AddMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID, roleID uint8, currentUserID uuid.UUID) error
	RemoveMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID, currentUserID uuid.UUID) error
}

type groupService struct {
	groupRepo        repositories.GroupRepository
	authorizeService AuthorizeService
}

func NewGroupService(groupRepo repositories.GroupRepository, authorizeService AuthorizeService) GroupService {
	return &groupService{
		groupRepo:        groupRepo,
		authorizeService: authorizeService,
	}
}

func (s *groupService) List(ctx context.Context, offset, limit int, search *string) ([]models.Group, error) {
	groups, err := s.groupRepo.List(ctx, offset, limit, search)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	return groups, nil
}

func (s *groupService) GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, autzerrors.NewGroupNotFoundError(id.String())
		}
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	return group, nil
}

func (s *groupService) Create(ctx context.Context, group *models.Group, ownerUserID uuid.UUID, currentUserID uuid.UUID) error {
	if group.ParentGroupID != nil {
		canCreate, err := s.authorizeService.CanCreateSubgroup(ctx, currentUserID, *group.ParentGroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return autzerrors.NewGroupNotFoundError(group.ParentGroupID.String())
			}
			return fmt.Errorf("failed to check subgroup creation permission: %w", err)
		}
		if !canCreate {
			return autzerrors.NewPermissionError("create subgroup")
		}
	} else {
		canCreate, err := s.authorizeService.CanCreateRootGroup(ctx, currentUserID)
		if err != nil {
			return fmt.Errorf("failed to check root group creation permission: %w", err)
		}
		if !canCreate {
			return autzerrors.NewPermissionError("create root group")
		}
	}

	if err := s.groupRepo.Create(ctx, group, ownerUserID); err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	return nil
}

func (s *groupService) Update(ctx context.Context, group *models.Group, ownerUserID *uuid.UUID, currentUserID uuid.UUID) error {
	canManage, err := s.authorizeService.CanManageGroup(ctx, currentUserID, group.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return autzerrors.NewGroupNotFoundError(group.ID.String())
		}
		return fmt.Errorf("failed to check group management permission: %w", err)
	}
	if !canManage {
		return autzerrors.NewPermissionError("update group")
	}

	if err := s.groupRepo.Update(ctx, group, ownerUserID); err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}
	return nil
}

func (s *groupService) Delete(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return autzerrors.NewGroupNotFoundError(id.String())
		}
		return fmt.Errorf("failed to get group: %w", err)
	}

	if len(group.SubGroups) > 0 {
		return autzerrors.NewValidationError("cannot delete group with subgroups")
	}

	canManage, err := s.authorizeService.CanManageGroup(ctx, currentUserID, id)
	if err != nil {
		return fmt.Errorf("failed to check group management permission: %w", err)
	}
	if !canManage {
		return autzerrors.NewPermissionError("delete group")
	}

	if err := s.groupRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}

func (s *groupService) GetMembers(ctx context.Context, groupID uuid.UUID, currentUserID uuid.UUID) ([]models.UserRole, error) {
	canView, err := s.authorizeService.CanViewGroupMembers(ctx, currentUserID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, autzerrors.NewGroupNotFoundError(groupID.String())
		}
		return nil, fmt.Errorf("failed to check member view permission: %w", err)
	}
	if !canView {
		return nil, autzerrors.NewPermissionError("view group members")
	}

	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	return members, nil
}

func (s *groupService) AddMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID, roleID uint8, currentUserID uuid.UUID) error {
	canManage, err := s.authorizeService.CanManageGroup(ctx, currentUserID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return autzerrors.NewGroupNotFoundError(groupID.String())
		}
		return fmt.Errorf("failed to check group management permission: %w", err)
	}
	if !canManage {
		return autzerrors.NewPermissionError("add group member")
	}

	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group members: %w", err)
	}

	for _, member := range members {
		if member.UserID == userID {
			return autzerrors.NewValidationError("user is already a member of this group")
		}
	}

	userRole := &models.UserRole{
		UserID:  userID,
		GroupID: groupID,
		RoleID:  roleID,
	}

	if err := s.groupRepo.AddMember(ctx, userRole); err != nil {
		return fmt.Errorf("failed to add member to group: %w", err)
	}
	return nil
}

func (s *groupService) RemoveMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID, currentUserID uuid.UUID) error {
	canManage, err := s.authorizeService.CanManageGroup(ctx, currentUserID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return autzerrors.NewGroupNotFoundError(groupID.String())
		}
		return fmt.Errorf("failed to check group management permission: %w", err)
	}
	if !canManage {
		return autzerrors.NewPermissionError("remove group member")
	}

	owner, err := s.groupRepo.GetOwner(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group owner: %w", err)
	}
	if owner.UserID == userID {
		return autzerrors.NewValidationError("cannot remove the group owner")
	}

	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group members: %w", err)
	}

	isMember := false
	for _, member := range members {
		if member.UserID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		return autzerrors.NewValidationError("user is not a member of this group")
	}

	if err := s.groupRepo.RemoveMember(ctx, groupID, userID); err != nil {
		return fmt.Errorf("failed to remove member from group: %w", err)
	}
	return nil
}

func (s *groupService) Count(ctx context.Context, search *string) (int64, error) {
	count, err := s.groupRepo.Count(ctx, search)
	if err != nil {
		return 0, fmt.Errorf("failed to count groups: %w", err)
	}
	return count, nil
}
