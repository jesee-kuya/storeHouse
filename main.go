package main

import (
	"storeHouse/database"
	hanlers "storeHouse/hanlers"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	database.ApplyMigrations(db)

	// Start the HTTP server
	hanlers.StartServer(db)
}
