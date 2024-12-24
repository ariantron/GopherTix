package database

import (
	"fmt"
	"gopher_tix/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func BuildDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		configs.DbHost,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbPort,
		configs.DbSslMode,
	)
}

func ConnectDB() *gorm.DB {
	dsn := BuildDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return db
}
