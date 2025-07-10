package routes

import (
	"github.com/gin-gonic/gin"
	"oe02_go_tam/database"
	"oe02_go_tam/handlers"
	"oe02_go_tam/middlewares"
	"oe02_go_tam/repositories"
	"oe02_go_tam/services"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authRepository := repositories.NewAuthRepository(database.DB)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)

	api := r.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.Use(middlewares.RequestLogger())
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/logout", authHandler.Logout)

	return r
}
