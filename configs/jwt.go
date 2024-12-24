package configs

import (
	"log"
	"os"
)

var SecretKey []byte

func init() {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("JWT_SECRET_KEY is not set in the .env file")
	}
	SecretKey = []byte(secretKey)
}
