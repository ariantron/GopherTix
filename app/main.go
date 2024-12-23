package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopher_tix/configs"
	autnmodels "gopher_tix/modules/authentication/models"
	autzmodels "gopher_tix/modules/authorization/models"
	profmodels "gopher_tix/modules/profile/models"
	tcktmodels "gopher_tix/modules/ticketing/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	loadEnv()
	dsn := buildDSN(configs.LoadDB())
	db := connectToDB(dsn)
	runMigrations(db)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func buildDSN(config *configs.DB) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort, config.DBSSLMode)
}

func connectToDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return db
}

func runMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&autnmodels.Login{},
		&autnmodels.User{},
		&autnmodels.Verify{},
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
