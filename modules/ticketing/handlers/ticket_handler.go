package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gopher_tix/modules/ticketing/models"
	"gopher_tix/modules/ticketing/requests"
	"gopher_tix/modules/ticketing/services"
	errs "gopher_tix/packages/common/errors"
	"gopher_tix/packages/utils"
)

type TicketHandler interface {
	GetByID(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	LeaveComment(ctx *fiber.Ctx) error
	RegisterRoutes(router fiber.Router)
}

type ticketHandler struct {
	service services.TicketService
}

func NewTicketHandler(service services.TicketService) TicketHandler {
	return &ticketHandler{
		service: service,
	}
}

func (h *ticketHandler) RegisterRoutes(router fiber.Router) {
	routes := router.Group("/tickets")
	routes.Get("/:id", h.GetByID)
	routes.Post("/", h.Create)
	routes.Post("/comment", h.LeaveComment)
}

func (h *ticketHandler) GetByID(ctx *fiber.Ctx) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	ticket, err := h.service.GetTicketByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(ticket)
}

func (h *ticketHandler) Create(ctx *fiber.Ctx) error {
	var req requests.CreateTicketRequest
	if err := errs.ParseAndValidateRequest(ctx, &req); err != nil {
		return err
	}

	ticket := &models.Ticket{
		Title:      req.Title,
		Text:       req.Text,
		GroupID:    req.GroupID,
		CreatedBy:  req.CreatedBy,
		AssignedTo: req.AssignedTo,
	}

	if err := h.service.CreateTicket(ctx.Context(), ticket); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(ticket)
}

func (h *ticketHandler) LeaveComment(ctx *fiber.Ctx) error {
	var req requests.LeaveCommentRequest
	if err := errs.ParseAndValidateRequest(ctx, &req); err != nil {
		return err
	}

	comment := &models.Comment{
		TicketID:    req.TicketID,
		CommentedBy: req.CommentedBy,
		Text:        req.Text,
	}

	if err := h.service.LeaveComment(ctx.Context(), comment); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(comment)
}
