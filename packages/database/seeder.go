package database

import (
	"gopher_tix/configs"
	autnseeders "gopher_tix/modules/authentication/seeders"
	autzseeders "gopher_tix/modules/authorization/seeders"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	autnseeders.CreateSuperAdmin(db)
	if configs.AppEnv == configs.DEV {
		autnseeders.UserSeeder(db)
	}
	autzseeders.RoleSeeder(db)
}
