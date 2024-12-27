package configs

import (
	"gopher_tix/packages/utils"
	"log"
	"os"
)

var SecretKey []byte

func init() {
	utils.LoadEnv()

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatalf("JWT_SECRET is not set in the .env file")
	}
	SecretKey = []byte(secretKey)
}
