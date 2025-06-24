package main

import (
	"oe02_go_tam/database"
	"oe02_go_tam/routes"
)

func main() {
	database.ConnectDB()

	r := routes.SetupRouter()
	r.Run(":8085")
}
