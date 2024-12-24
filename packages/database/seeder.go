package database

import (
	autnseeders "gopher_tix/modules/authentication/seeders"
	autzseeders "gopher_tix/modules/authorization/seeders"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	autnseeders.CreateSuperAdmin(db)
	autnseeders.UserSeeder(db)
	autzseeders.RoleSeeder(db)
}
