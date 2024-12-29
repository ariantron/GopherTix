package services

import (
	"context"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/constants"
	autzservices "gopher_tix/modules/authorization/services"
	"gopher_tix/modules/ticketing/models"
	"gopher_tix/modules/ticketing/repositories"
	errs "gopher_tix/packages/common/errors"
)

type TicketService interface {
	GetTicketByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error)
	CreateTicket(ctx context.Context, ticket *models.Ticket) error
	LeaveComment(ctx context.Context, comment *models.Comment) error
}

type ticketService struct {
	ticketRepo    repositories.TicketRepository
	authorizeServ autzservices.AuthorizeService
}

func NewTicketService(
	ticketRepo repositories.TicketRepository,
	authorizeServ autzservices.AuthorizeService,
) TicketService {
	return &ticketService{
		ticketRepo:    ticketRepo,
		authorizeServ: authorizeServ,
	}
}

func (s *ticketService) GetTicketByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	return s.ticketRepo.GetTicketByID(ctx, id)
}

func (s *ticketService) CreateTicket(ctx context.Context, ticket *models.Ticket) error {
	canCreate, err := s.authorizeServ.CanCreateTicket(ctx, ticket.CreatedBy, ticket.GroupID)
	if err != nil {
		return err
	}
	if !canCreate {
		return errs.NewPermissionDeniedError("Create Ticket")
	}

	isMember, err := s.authorizeServ.HasRole(ctx, ticket.AssignedTo, ticket.GroupID, constants.MemberRoleID)
	if err != nil {
		return err
	}
	if !isMember {
		return errs.NewIncorrectParameter("Recipient must be a member of the group")
	}

	return s.ticketRepo.CreateTicket(ctx, ticket)
}

func (s *ticketService) LeaveComment(ctx context.Context, comment *models.Comment) error {
	ticket, err := s.ticketRepo.GetTicketByID(ctx, comment.TicketID)
	if err != nil {
		return err
	}

	canComment, err := s.authorizeServ.CanLeaveComment(
		ctx,
		comment.CommentedBy,
		ticket.CreatedBy,
		ticket.AssignedTo,
		ticket.GroupID,
	)
	if err != nil {
		return err
	}
	if !canComment {
		return errs.NewPermissionDeniedError("Leave Comment")
	}

	return s.ticketRepo.LeaveComment(ctx, comment)
}
