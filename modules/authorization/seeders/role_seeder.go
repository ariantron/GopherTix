package seeders

import (
	"gopher_tix/modules/authorization/constants"
	"gopher_tix/modules/authorization/models"
	"gorm.io/gorm"
	"log"
)

func RoleSeeder(db *gorm.DB) {
	var count int64
	db.Model(&models.Role{}).Count(&count)

	if count == 0 {
		roleNames := [...]string{
			constants.AdminRole,
			constants.OwnerRole,
			constants.MemberRole,
		}

		for _, roleName := range roleNames {
			role := models.Role{
				Name: roleName,
			}
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to seed role: %v", err)
			}
		}
		log.Println("Role seeding completed successfully")
	} else {
		log.Println("Roles already seeded, skipping...")
	}
}
