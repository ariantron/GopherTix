package configs

import (
	"log"
	"os"
)

var (
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
)

func init() {
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	DBSSLMode = os.Getenv("DB_SSLMODE")

	if DBHost == "" || DBPort == "" || DBUser == "" || DBPassword == "" || DBName == "" || DBSSLMode == "" {
		log.Fatalf("Missing required environment variable(s) for DB configuration")
	}
}
