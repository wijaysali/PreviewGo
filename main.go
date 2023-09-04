package main

import (
	"photo-app/database"
	"photo-app/router"
)

func main() {
	database.StartDB()
	r := router.StartApp()
	r.Run(":8080")
}
