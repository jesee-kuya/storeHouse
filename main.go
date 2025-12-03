package main

import (
	. "storeHouse/database"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db := ConnectDB()
	defer db.Close()

	ApplyMigrations(db)

}
