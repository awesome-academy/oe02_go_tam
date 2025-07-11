package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
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

	tourRepo := repositories.NewTourRepository(database.DB)
	tourService := services.NewTourService(tourRepo)
	tourHandler := handlers.NewTourHandler(tourService)

	reviewRepo := repositories.NewReviewRepository(database.DB)
	reviewService := services.NewReviewService(reviewRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	api := r.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.Use(middlewares.RequestLogger())
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/logout", authHandler.Logout)
	authGroup.GET("/google", func(c *gin.Context) {
		c.Request.URL.RawQuery = "provider=google"
		gothic.BeginAuthHandler(c.Writer, c.Request)
	})
	authGroup.GET("/google/callback", authHandler.GoogleCallback)

	userGroup := api.Group("/users")
	userGroup.Use(middlewares.RequestLogger())
	userGroup.Use(middlewares.AuthMiddleware())
	userGroup.GET("/me", userHandler.GetProfile)
	userGroup.PATCH("/me", userHandler.UpdateProfile)

	tourGroup := api.Group("/tours")
	tourGroup.Use(middlewares.RequestLogger())
	tourGroup.GET("/", tourHandler.ListTours)
	tourGroup.GET("/:id", tourHandler.GetTourDetail)
	tourGroup.GET("/:id/reviews", reviewHandler.GetReviews)

	return r
}
