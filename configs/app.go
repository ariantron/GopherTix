package configs

import (
	"log"
	"os"
)

var (
	AppName string
	AppPort string
)

func init() {
	AppName = os.Getenv("APP_NAME")
	AppPort = os.Getenv("APP_PORT")

	if AppName == "" || AppPort == "" {
		log.Fatalf("Missing required environment variable(s) for APP configuration")
	}
}
