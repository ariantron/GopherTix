package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/requests"
	"gopher_tix/modules/authentication/services"
	errs "gopher_tix/packages/common/errors"
	"net"
)

type LoginHandler interface {
	Login(ctx *fiber.Ctx) error
	RegisterRoutes(router fiber.Router)
}

type loginHandler struct {
	loginService services.LoginService
}

func NewLoginHandler(loginService services.LoginService) LoginHandler {
	return &loginHandler{
		loginService: loginService,
	}
}

func (h *loginHandler) RegisterRoutes(router fiber.Router) {
	routes := router.Group("/auth")
	routes.Post("/login", h.Login)
}

func (h *loginHandler) Login(ctx *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := errs.ParseAndValidateRequest(ctx, req); err != nil {
		return err
	}

	ip := net.ParseIP(ctx.IP())
	var userID uuid.UUID
	var userEmail string
	var token string

	user, token, err := h.loginService.ValidateUserCredentials(ctx.Context(), req.Email, req.Password)
	if user != nil {
		userID = user.ID
		userEmail = user.Email
	}

	loginRecord := &models.Login{
		UserID:  userID,
		IP:      ip,
		Succeed: err == nil,
	}

	if err := h.loginService.CreateLoginRecord(ctx.Context(), loginRecord); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    userID,
			"email": userEmail,
		},
	})
}
