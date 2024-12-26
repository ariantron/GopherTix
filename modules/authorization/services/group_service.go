package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/constants"
	"gopher_tix/modules/authorization/models"
	"gopher_tix/modules/authorization/repositories"
)

type GroupServiceInterface interface {
	List(ctx context.Context, offset, limit int, search *string) ([]models.Group, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error)
	Create(ctx context.Context, group *models.Group, ownerUserID uuid.UUID, currentUserID uuid.UUID) error
	Update(ctx context.Context, group *models.Group, ownerUserID *uuid.UUID, currentUserID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error
	GetMembers(ctx context.Context, groupID uuid.UUID, currentUserID uuid.UUID) ([]models.UserRole, error)
	Count(ctx context.Context, search *string) (int64, error)
}

type GroupService struct {
	groupRepo     *repositories.GroupRepository
	authorizeRepo *repositories.AuthorizeRepository
}

func NewGroupService(groupRepo *repositories.GroupRepository, authorizeRepo *repositories.AuthorizeRepository) *GroupService {
	return &GroupService{
		groupRepo:     groupRepo,
		authorizeRepo: authorizeRepo,
	}
}

func (s *GroupService) List(ctx context.Context, offset, limit int, search *string) ([]models.Group, error) {
	return s.groupRepo.List(ctx, offset, limit, search)
}

func (s *GroupService) GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	return s.groupRepo.GetByID(ctx, id)
}

func (s *GroupService) Create(ctx context.Context, group *models.Group, ownerUserID uuid.UUID, currentUserID uuid.UUID) error {
	if group.ParentGroupID != nil {
		isAdmin, err := s.authorizeRepo.IsAdmin(ctx, currentUserID)
		if err != nil {
			return err
		}
		if !isAdmin {
			hasOwnerRole, err := s.authorizeRepo.HasRole(ctx, currentUserID, *group.ParentGroupID, constants.OwnerRoleID)
			if err != nil {
				return err
			}
			if !hasOwnerRole {
				return errors.New("unauthorized: only admin or parent group owner can create subgroups")
			}
		}
	} else {
		isAdmin, err := s.authorizeRepo.IsAdmin(ctx, currentUserID)
		if err != nil {
			return err
		}
		if !isAdmin {
			return errors.New("unauthorized: only admin can create root groups")
		}
	}
	return s.groupRepo.Create(ctx, group, ownerUserID)
}

func (s *GroupService) Update(ctx context.Context, group *models.Group, ownerUserID *uuid.UUID, currentUserID uuid.UUID) error {
	isAdmin, err := s.authorizeRepo.IsAdmin(ctx, currentUserID)
	if err != nil {
		return err
	}
	if !isAdmin {
		hasOwnerRole, err := s.authorizeRepo.HasRole(ctx, currentUserID, group.ID, constants.OwnerRoleID)
		if err != nil {
			return err
		}
		if !hasOwnerRole {
			return errors.New("unauthorized: only admin or group owner can update group")
		}
	}
	return s.groupRepo.Update(ctx, group, ownerUserID)
}

func (s *GroupService) Delete(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if len(group.SubGroups) > 0 {
		return errors.New("cannot delete group with subgroups")
	}
	isAdmin, err := s.authorizeRepo.IsAdmin(ctx, currentUserID)
	if err != nil {
		return err
	}
	if !isAdmin {
		hasOwnerRole, err := s.authorizeRepo.HasRole(ctx, currentUserID, id, constants.OwnerRoleID)
		if err != nil {
			return err
		}
		if !hasOwnerRole {
			return errors.New("unauthorized: only admin or group owner can delete group")
		}
	}
	return s.groupRepo.Delete(ctx, id)
}

func (s *GroupService) GetMembers(ctx context.Context, groupID uuid.UUID, currentUserID uuid.UUID) ([]models.UserRole, error) {
	isAdmin, err := s.authorizeRepo.IsAdmin(ctx, currentUserID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		hasAnyRole, err := s.authorizeRepo.HasRole(ctx, currentUserID, groupID, constants.MemberRoleID)
		if err != nil {
			return nil, err
		}
		if !hasAnyRole {
			return nil, errors.New("unauthorized: only admin or group members can view members")
		}
	}
	return s.groupRepo.GetMembers(ctx, groupID)
}

func (s *GroupService) Count(ctx context.Context, search *string) (int64, error) {
	return s.groupRepo.Count(ctx, search)
}
