package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"oe02_go_tam/auth"
	"oe02_go_tam/constant"
	"oe02_go_tam/database"
	_ "oe02_go_tam/docs"
	"oe02_go_tam/routes"
	"oe02_go_tam/utils"
)

func main() {
	database.ConnectDB()

	if err := utils.InitJWTSecret(); err != nil {
		log.Fatal("Failed to initialize JWT secret:", err)
		return
	}

	if err := auth.InitGoogleProvider(); err != nil {
		log.Fatalf("Failed to initialize Google OAuth: %v", err)
	}

	constant.LoadI18n("en")

	r := routes.SetupRouter()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8085")
}
