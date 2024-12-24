package configs

import (
	"gopher_tix/packages/initializer"
	"log"
	"os"
)

var SecretKey []byte

func init() {
	initializer.LoadEnv()

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatalf("JWT_SECRET is not set in the .env file")
	}
	SecretKey = []byte(secretKey)
}
