package configs

import (
	"log"
	"os"
)

var (
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbSslMode  string
)

func init() {
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")
	DbSslMode = os.Getenv("DB_SSLMODE")

	if DbHost == "" || DbPort == "" || DbUser == "" || DbPassword == "" || DbName == "" || DbSslMode == "" {
		log.Fatalf("Missing required environment variable(s) for DB configuration")
	}
}
