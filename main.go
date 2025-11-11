package main

import (
	. "storeHouse/database"
)

func main() {
	db := ConnectDB()
	defer db.Close()
}
