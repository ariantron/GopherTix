package seeders

import (
	"github.com/go-faker/faker/v4"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/packages/utils"
	"gorm.io/gorm"
	"log"
)

func CreateSuperAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("email = ?", "superadmin@gophertix.com").Count(&count)

	if count > 0 {
		log.Println("SuperAdmin already exists. Skipping creation.")
		return
	}

	hashedPassword, err := utils.HashPassword("Abc123")
	if err != nil {
		log.Printf("Failed to hash password")
		return
	}

	user := models.User{
		Email:    "superadmin@gophertix.com",
		Password: hashedPassword,
	}
	if err := db.Create(&user).Error; err != nil {
		log.Printf("Failed to create SuperAdmin")
	} else {
		log.Println("SuperAdmin created successfully.")
	}
}

func UserSeeder(db *gorm.DB) {
	const numberOfUsers = 11

	var count int64
	db.Model(&models.User{}).Count(&count)

	if count >= numberOfUsers {
		log.Printf("Database already contains %d or more users. Skipping user seeding.", count)
		return
	}

	for i := 0; i < numberOfUsers-int(count); i++ {
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

	log.Printf("%d users seeded successfully.", numberOfUsers-int(count))
}
