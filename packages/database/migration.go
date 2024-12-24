package database

import (
	"fmt"
	autnmodels "gopher_tix/modules/authentication/models"
	autzmodels "gopher_tix/modules/authorization/models"
	profmodels "gopher_tix/modules/profile/models"
	tcktmodels "gopher_tix/modules/ticketing/models"
	"gorm.io/gorm"
	"log"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&autnmodels.Login{},
		&autnmodels.User{},
		&autzmodels.Group{},
		&autzmodels.Permission{},
		&autzmodels.Role{},
		&autzmodels.RolePermission{},
		&autzmodels.UserPermission{},
		&autzmodels.UserRole{},
		&profmodels.Profile{},
		&profmodels.ProfileField{},
		&profmodels.Field{},
		&tcktmodels.Ticket{},
		&tcktmodels.Comment{},
	)
	if err != nil {
		log.Fatalf("Error migrating the models: %v", err)
	}
	fmt.Println("Database migrated successfully")
}
