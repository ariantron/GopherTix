package services

import (
	"context"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/models"
	"gopher_tix/modules/authorization/repositories"
	errs "gopher_tix/packages/common/errors"
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
	return s.groupRepo.List(ctx, offset, limit, search)
}

func (s *groupService) GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	return s.groupRepo.GetByID(ctx, id)
}

func (s *groupService) Create(ctx context.Context, group *models.Group, ownerUserID uuid.UUID, currentUserID uuid.UUID) error {
	if group.ParentGroupID != nil {
		canCreate, err := s.authorizeService.CanCreateSubgroup(ctx, currentUserID, *group.ParentGroupID)
		if err != nil {
			return err
		}
		if !canCreate {
			return errs.NewPermissionDeniedError("Create Subgroup")
		}
	} else {
		canCreate, err := s.authorizeService.CanCreateRootGroup(ctx, currentUserID)
		if err != nil {
			return err
		}
		if !canCreate {
			return errs.NewPermissionDeniedError("Create Root Group")
		}
	}

	return s.groupRepo.Create(ctx, group, ownerUserID)
}

func (s *groupService) Update(ctx context.Context, group *models.Group, ownerUserID *uuid.UUID, currentUserID uuid.UUID) error {
	canManage, err := s.authorizeService.CanManageGroup(ctx, currentUserID, group.ID)
	if err != nil {
		return err
	}
	if !canManage {
		return errs.NewPermissionDeniedError("Update Group")
	}

	return s.groupRepo.Update(ctx, group, ownerUserID)
}

func (s *groupService) Delete(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) error {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if len(group.SubGroups) > 0 {
		return errs.NewIncorrectParameter("Cannot delete group with subgroups")
	}

	canManage, err := s.authorizeService.CanManageGroup(ctx, currentUserID, id)
	if err != nil {
		return err
	}
	if !canManage {
		return errs.NewPermissionDeniedError("Delete Group")
	}

	return s.groupRepo.Delete(ctx, id)
}

func (s *groupService) GetMembers(ctx context.Context, groupID uuid.UUID, currentUserID uuid.UUID) ([]models.UserRole, error) {
	canView, err := s.authorizeService.CanViewGroupMembers(ctx, currentUserID, groupID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, errs.NewPermissionDeniedError("View Group Members")
	}

	return s.groupRepo.GetMembers(ctx, groupID)
}

func (s *groupService) AddMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID, roleID uint8, currentUserID uuid.UUID) error {
	canManage, err := s.authorizeService.CanManageGroup(ctx, currentUserID, groupID)
	if err != nil {
		return err
	}
	if !canManage {
		return errs.NewPermissionDeniedError("Add Group Member")
	}

	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return err
	}

	for _, member := range members {
		if member.UserID == userID {
			return errs.NewIncorrectParameter("User is already a member of this group")
		}
	}

	userRole := &models.UserRole{
		UserID:  userID,
		GroupID: groupID,
		RoleID:  roleID,
	}

	return s.groupRepo.AddMember(ctx, userRole)
}

func (s *groupService) RemoveMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID, currentUserID uuid.UUID) error {
	canManage, err := s.authorizeService.CanManageGroup(ctx, currentUserID, groupID)
	if err != nil {
		return err
	}
	if !canManage {
		return errs.NewPermissionDeniedError("Remove Group Member")
	}

	owner, err := s.groupRepo.GetOwner(ctx, groupID)
	if err != nil {
		return err
	}
	if owner.UserID == userID {
		return errs.NewIncorrectParameter("Cannot remove the group owner")
	}

	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return err
	}

	isMember := false
	for _, member := range members {
		if member.UserID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		return errs.NewIncorrectParameter("User is not a member of this group")
	}

	return s.groupRepo.RemoveMember(ctx, groupID, userID)
}

func (s *groupService) Count(ctx context.Context, search *string) (int64, error) {
	return s.groupRepo.Count(ctx, search)
}
