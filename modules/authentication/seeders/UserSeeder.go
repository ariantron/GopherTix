package seeders

import (
	"github.com/go-faker/faker/v4"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/packages/common/utils"
	"gorm.io/gorm"
	"log"
)

func UserSeeder(db *gorm.DB) {
	const numberOfUsers = 10

	for i := 0; i < numberOfUsers; i++ {
		hashedPassword, err := utils.HashPassword(faker.Password())
		if err != nil {
			log.Printf("Failed to hash password for user %d: %v", i+1, err)
			continue
		}
		user := models.User{
			Email:    faker.Email(),
			Password: hashedPassword,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %d: %v", i+1, err)
		}
	}

	log.Printf("%d users seeded successfully.", numberOfUsers)
}
