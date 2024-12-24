package main

import (
	"gopher_tix/packages/database"
	"gopher_tix/server"
)

func main() {
	db := database.ConnectDB()
	database.Migrate(db)
	database.Seed(db)
	server.Serve(db)
}
