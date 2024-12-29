package router

import (
	"github.com/gofiber/fiber/v2"
	autzrepos "gopher_tix/modules/authorization/repositories"
	autzservices "gopher_tix/modules/authorization/services"
	tckthandlers "gopher_tix/modules/ticketing/handlers"
	tcktrepos "gopher_tix/modules/ticketing/repositories"
	tcktservices "gopher_tix/modules/ticketing/services"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router fiber.Router) {
	ticketRepo := tcktrepos.NewTicketRepository(db)
	authorizeRepo := autzrepos.NewAuthorizeRepository(db)
	groupRepo := autzrepos.NewGroupRepository(db)

	authorizeService := autzservices.NewAuthorizeService(authorizeRepo, groupRepo)
	ticketService := tcktservices.NewTicketService(ticketRepo, authorizeService)

	ticketHandler := tckthandlers.NewTicketHandler(ticketService)

	ticketHandler.RegisterRoutes(router)
}
