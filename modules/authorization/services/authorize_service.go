package services

import (
	"context"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/constants"
	"gopher_tix/modules/authorization/repositories"
)

type AuthorizeService interface {
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	HasRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) (bool, error)
	CanCreateSubgroup(ctx context.Context, userID uuid.UUID, parentGroupID uuid.UUID) (bool, error)
	CanCreateRootGroup(ctx context.Context, userID uuid.UUID) (bool, error)
	CanManageGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error)
	CanViewGroupMembers(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error)
	CanCreateTicket(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error)
	CanLeaveComment(ctx context.Context, userID uuid.UUID, ticketCreatorID uuid.UUID, ticketRecipientID uuid.UUID, groupID uuid.UUID) (bool, error)
}

type authorizeService struct {
	authorizeRepo repositories.AuthorizeRepository
	groupRepo     repositories.GroupRepository
}

func NewAuthorizeService(authorizeRepo repositories.AuthorizeRepository, groupRepo repositories.GroupRepository) AuthorizeService {
	return &authorizeService{
		authorizeRepo: authorizeRepo,
		groupRepo:     groupRepo,
	}
}

func (s *authorizeService) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	return s.authorizeRepo.IsAdmin(ctx, userID)
}

func (s *authorizeService) HasRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) (bool, error) {
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	return s.authorizeRepo.HasRole(ctx, userID, groupID, roleID)
}

func (s *authorizeService) CanCreateSubgroup(ctx context.Context, userID uuid.UUID, parentGroupID uuid.UUID) (bool, error) {
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	return s.HasRole(ctx, userID, parentGroupID, constants.OwnerRoleID)
}

func (s *authorizeService) CanCreateRootGroup(ctx context.Context, userID uuid.UUID) (bool, error) {
	return s.IsAdmin(ctx, userID)
}

func (s *authorizeService) CanManageGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	return s.HasRole(ctx, userID, groupID, constants.OwnerRoleID)
}

func (s *authorizeService) CanViewGroupMembers(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	return s.HasRole(ctx, userID, groupID, constants.MemberRoleID)
}

func (s *authorizeService) CanCreateTicket(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}

	isOwner, err := s.IsOwnerOfGroupOrParents(ctx, userID, groupID)
	if err != nil {
		return false, err
	}
	if isOwner {
		return true, nil
	}

	return s.HasRole(ctx, userID, groupID, constants.MemberRoleID)
}

func (s *authorizeService) CanLeaveComment(ctx context.Context, userID uuid.UUID, ticketCreatorID uuid.UUID, ticketRecipientID uuid.UUID, groupID uuid.UUID) (bool, error) {
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}

	isOwner, err := s.IsOwnerOfGroupOrParents(ctx, userID, groupID)
	if err != nil {
		return false, err
	}
	if isOwner {
		return true, nil
	}

	if userID == ticketCreatorID || userID == ticketRecipientID {
		return true, nil
	}

	return false, nil
}

func (s *authorizeService) IsOwnerOfGroupOrParents(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	currentGroupID := groupID
	for {
		isOwner, err := s.HasRole(ctx, userID, currentGroupID, constants.OwnerRoleID)
		if err != nil {
			return false, err
		}
		if isOwner {
			return true, nil
		}

		group, err := s.groupRepo.GetByID(ctx, currentGroupID)
		if err != nil {
			return false, err
		}

		if group.ParentGroupID == nil {
			break
		}

		currentGroupID = *group.ParentGroupID
	}

	return false, nil
}
