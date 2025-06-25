package routes

import (
	"github.com/gin-gonic/gin"
	"oe02_go_tam/handlers"
	"oe02_go_tam/repositories"
	"oe02_go_tam/services"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authRepository := repositories.NewAuthRepository()
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)

	api := r.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)

	return r
}
