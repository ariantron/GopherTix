package main

import (
	"gopher_tix/configs"
	"gopher_tix/packages/database"
	"gopher_tix/server"
)

func main() {
	db := database.ConnectDB()
	database.Migrate(db)
	if configs.AppEnv == configs.DEV {
		database.Seed(db)
	}
	server.Serve(db)
}
