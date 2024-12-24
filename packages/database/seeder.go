package database

import (
	"gopher_tix/modules/authentication/seeders"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seeders.UserSeeder(db)
}
