package main

import (
	"gopher_tix/packages/database"
	"gopher_tix/packages/init"
)

func main() {
	init.LoadEnv()
	db := database.ConnectDB()
	database.RunMigrations(db)
	serve(db)
}
