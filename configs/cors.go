package configs

import (
	"gopher_tix/packages/utils"
	"log"
	"os"
)

var (
	CorsAllowOrigins string
	CorsAllowHeaders string
	CorsAllowMethods string
)

func init() {
	utils.LoadEnv()

	CorsAllowOrigins = os.Getenv("CORS_ALLOW_ORIGINS")
	CorsAllowHeaders = os.Getenv("CORS_ALLOW_HEADERS")
	CorsAllowMethods = os.Getenv("CORS_ALLOW_METHODS")

	if CorsAllowOrigins == "" || CorsAllowHeaders == "" || CorsAllowMethods == "" {
		log.Fatalf("Missing required cors environment variable(s) for APP configuration")
	}
}
