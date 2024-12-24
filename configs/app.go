package configs

import (
	"gopher_tix/packages/initializer"
	"log"
	"os"
)

var (
	AppName string
	AppEnv  string
	AppPort string
)

const (
	DEV  = "dev"
	PROD = "prod"
)

func init() {
	initializer.LoadEnv()

	AppName = os.Getenv("APP_NAME")
	AppEnv = os.Getenv("APP_ENV")
	AppPort = os.Getenv("APP_PORT")

	if AppName == "" || AppPort == "" || AppEnv == "" {
		log.Fatalf("Missing required app environment variable(s) for APP configuration")
	}

	if AppEnv != PROD && AppEnv != DEV {
		log.Fatalf("Invalid APP_ENV value: %s. Expected 'prod' or 'dev'.", AppEnv)
	}

	if AppEnv == PROD {
		log.Println("Running in Production mode.")
	} else if AppEnv == DEV {
		log.Println("Running in Development mode.")
	}
}
