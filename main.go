package main

import (
	"log"
	"oe02_go_tam/constant"
	"oe02_go_tam/database"
	"oe02_go_tam/routes"
	"oe02_go_tam/utils"
)

func main() {
	database.ConnectDB()

	if err := utils.InitJWTSecret(); err != nil {
		log.Fatal("Failed to initialize JWT secret:", err)
		return
	}

	constant.LoadI18n("en")

	r := routes.SetupRouter()
	r.Run(":8085")
}
