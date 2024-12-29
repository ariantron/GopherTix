package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gopher_tix/modules/ticketing/models"
	errs "gopher_tix/packages/common/errors"
	"gorm.io/gorm"
)

type TicketRepository interface {
	GetTicketByID(ctx context.Context, ticketID uuid.UUID) (*models.Ticket, error)
	CreateTicket(ctx context.Context, ticket *models.Ticket) error
	LeaveComment(ctx context.Context, comment *models.Comment) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{
		db: db,
	}
}

func (r *ticketRepository) GetTicketByID(ctx context.Context, ticketID uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	result := r.db.WithContext(ctx).First(&ticket, "id = ?", ticketID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("Ticket")
		}
		return nil, result.Error
	}
	return &ticket, nil
}

func (r *ticketRepository) CreateTicket(ctx context.Context, ticket *models.Ticket) error {
	err := r.db.WithContext(ctx).Create(ticket).Error
	if err != nil {
		return errs.NewInternalServerError("Failed to create ticket")
	}
	return nil
}

func (r *ticketRepository) LeaveComment(ctx context.Context, comment *models.Comment) error {
	var ticket models.Ticket
	result := r.db.WithContext(ctx).First(&ticket, "id = ?", comment.TicketID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errs.NewNotFoundError("Ticket")
		}
		return result.Error
	}

	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		return errs.NewInternalServerError("Failed to leave comment")
	}

	return nil
}
