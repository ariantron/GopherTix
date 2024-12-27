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
}

type authorizeService struct {
	repo repositories.AuthorizeRepository
}

func NewAuthorizeService(repo repositories.AuthorizeRepository) AuthorizeService {
	return &authorizeService{repo: repo}
}

func (s *authorizeService) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	return s.repo.IsAdmin(ctx, userID)
}

func (s *authorizeService) HasRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) (bool, error) {
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	return s.repo.HasRole(ctx, userID, groupID, roleID)
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
