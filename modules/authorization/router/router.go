package router

import (
	"github.com/gofiber/fiber/v2"
	"gopher_tix/modules/authorization/handlers"
	"gopher_tix/modules/authorization/repositories"
	"gopher_tix/modules/authorization/services"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router fiber.Router) {
	authorizeRepo := repositories.NewAuthorizeRepository(db)
	groupRepo := repositories.NewGroupRepository(db)

	authorizeService := services.NewAuthorizeService(authorizeRepo, groupRepo)
	groupService := services.NewGroupService(groupRepo, authorizeService)

	groupHandler := handlers.NewGroupHandler(groupService)

	groupHandler.RegisterRoutes(router)
}
