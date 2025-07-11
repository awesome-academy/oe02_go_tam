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
	userRepository := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	api := r.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.Use(middlewares.RequestLogger())
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/logout", authHandler.Logout)

	userGroup := api.Group("/users")
	userGroup.Use(middlewares.RequestLogger())
	userGroup.Use(middlewares.AuthMiddleware())
	userGroup.GET("/me", userHandler.GetProfile)
	userGroup.PATCH("/me", userHandler.UpdateProfile)

	return r
}
