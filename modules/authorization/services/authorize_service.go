package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/constants"
	"gopher_tix/modules/authorization/models"
	"gopher_tix/modules/authorization/repositories"
)

type AuthorizeServiceInterface interface {
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	HasRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) (bool, error)
	AssignRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) error
	UnassignRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) error
}

type AuthorizeService struct {
	repo *repositories.AuthorizeRepository
}

func NewAuthorizeService(repo *repositories.AuthorizeRepository) *AuthorizeService {
	return &AuthorizeService{repo: repo}
}

func (s *AuthorizeService) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	return s.repo.IsAdmin(ctx, userID)
}

func (s *AuthorizeService) HasRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) (bool, error) {
	isAdmin, err := s.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	return s.repo.HasRole(ctx, userID, groupID, roleID)
}

func (s *AuthorizeService) AssignRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) error {
	if roleID == constants.AdminRoleID && groupID != uuid.Nil {
		return errors.New("admin role can only be assigned at system level")
	}
	if roleID == constants.OwnerRoleID {
		hasOwner, err := s.repo.HasRole(ctx, uuid.Nil, groupID, constants.OwnerRoleID)
		if err != nil {
			return err
		}
		if hasOwner {
			return errors.New("group already has an owner")
		}
	}
	userRole := &models.UserRole{
		UserID:  userID,
		GroupID: groupID,
		RoleID:  roleID,
	}
	return s.repo.AssignRole(ctx, userRole)
}

func (s *AuthorizeService) UnassignRole(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, roleID uint8) error {
	if roleID == constants.OwnerRoleID {
		hasOwner, err := s.repo.HasRole(ctx, uuid.Nil, groupID, constants.OwnerRoleID)
		if err != nil {
			return err
		}
		if !hasOwner {
			return errors.New("cannot remove the only owner of a group")
		}
	}
	userRole := &models.UserRole{
		UserID:  userID,
		GroupID: groupID,
		RoleID:  roleID,
	}
	return s.repo.UnassignRole(ctx, userRole)
}
