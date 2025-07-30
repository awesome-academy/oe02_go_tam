package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"oe02_go_tam/config"
	"oe02_go_tam/database"
	"oe02_go_tam/handlers"
	"oe02_go_tam/middlewares"
	"oe02_go_tam/repositories"
	"oe02_go_tam/services"
	"oe02_go_tam/utils"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//r.LoadHTMLGlob("templates/*.html")
	utils.LoadTemplates(r)

	authRepository := repositories.NewAuthRepository(database.DB)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)
	adminAuthHandler := handlers.NewAdminAuthHandler(authService)
	userRepository := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	tourRepo := repositories.NewTourRepository(database.DB)
	tourService := services.NewTourService(tourRepo)
	tourHandler := handlers.NewTourHandler(tourService)

	reviewRepo := repositories.NewReviewRepository(database.DB)
	reviewService := services.NewReviewService(reviewRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	commentRepo := repositories.NewCommentRepository(database.DB)
	commentService := services.NewCommentService(commentRepo, reviewRepo)
	commentHandler := handlers.NewCommentHandler(commentService)

	likeRepo := repositories.NewLikeRepository(database.DB)
	likeService := services.NewLikeService(likeRepo, reviewRepo)
	likeHandler := handlers.NewLikeHandler(likeService)

	bookingRepo := repositories.NewBookingRepository(database.DB)
	bookingService := services.NewBookingService(bookingRepo, tourRepo)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	transactionRepo := repositories.NewTransactionRepository(database.DB)
	vnpayConfig := config.GetVnpayConfig()
	vnpayService := services.NewVnpayService(bookingRepo, tourRepo, transactionRepo, vnpayConfig)
	vnpayHandler := handlers.NewVnpayHandler(vnpayService)

	adminUsersService := services.NewAdminUsersService(userRepository)
	adminUserHandler := handlers.NewAdminUsersHandler(adminUsersService)
	adminTourService := services.NewAdminTourService(tourRepo)
	adminTourHandler := handlers.NewAdminTourHandler(adminTourService)
	adminBookingService := services.NewAdminBookingService(bookingRepo)
	adminBookingHandler := handlers.NewAdminBookingHandler(adminBookingService)
	adminReviewService := services.NewAdminReviewService(reviewRepo)
	adminReviewHandler := handlers.NewAdminReviewHandler(adminReviewService)
	adminRevenueService := services.NewAdminRevenueService(bookingRepo)
	adminRevenueHandler := handlers.NewAdminRevenueHandler(adminRevenueService)

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

	adminGroup := r.Group("/admin")
	adminGroup.GET("/login", adminAuthHandler.ShowLoginForm)
	adminGroup.POST("/login", adminAuthHandler.HandleLogin)
	adminGroup.POST("/logout", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminAuthHandler.HandleLogout)
	adminGroup.GET("/users", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminUserHandler.ListUsers)
	adminGroup.GET("/users/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminUserHandler.ViewUser)
	adminGroup.GET("/users/:id/ban", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminUserHandler.ToggleBanUser)
	adminGroup.GET("/users/:id/delete", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminUserHandler.DeleteUser)
	adminGroup.GET("/tours", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminTourHandler.ListTours)
	adminGroup.GET("/tours/new", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminTourHandler.ShowCreateForm)
	adminGroup.POST("/tours", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminTourHandler.HandleCreate)
	adminGroup.GET("/tours/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminTourHandler.ViewTour)
	adminGroup.GET("/tours/:id/edit", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminTourHandler.ShowEditForm)
	adminGroup.POST("/tours/:id/edit", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminTourHandler.HandleUpdate)
	adminGroup.GET("/tours/:id/delete", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminTourHandler.DeleteTour)
	adminGroup.GET("/bookings", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminBookingHandler.ListBookings)
	adminGroup.GET("/bookings/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminBookingHandler.ShowBookingDetail)
	adminGroup.GET("/bookings/:id/delete", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminBookingHandler.DeleteBooking)
	adminGroup.GET("/bookings/:id/cancel", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminBookingHandler.CancelBooking)
	adminGroup.GET("/reviews", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminReviewHandler.ListReviews)
	adminGroup.GET("/reviews/:id", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminReviewHandler.ShowReviewDetail)
	adminGroup.GET("/reviews/:id/delete", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminReviewHandler.DeleteReview)
	adminGroup.GET("/revenues", middlewares.AuthMiddleware(), middlewares.RequireRole("admin"), adminRevenueHandler.ListRevenue)

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

	reviewGroup := api.Group("/reviews")
	reviewGroup.Use(middlewares.AuthMiddleware())
	reviewGroup.POST("", reviewHandler.CreateReview)
	reviewGroup.GET("/:id", reviewHandler.GetOwnReview)
	reviewGroup.PUT("/:id", reviewHandler.UpdateReview)
	reviewGroup.DELETE("/:id", reviewHandler.DeleteReview)

	commentGroup := api.Group("/comments")
	commentGroup.Use(middlewares.AuthMiddleware())
	commentGroup.POST("", commentHandler.CreateComment)

	likeGroup := api.Group("/likes")
	likeGroup.Use(middlewares.AuthMiddleware())
	likeGroup.POST("", likeHandler.LikeReview)

	bookingGroup := api.Group("/bookings")
	bookingGroup.Use(middlewares.RequestLogger())
	bookingGroup.Use(middlewares.AuthMiddleware())
	bookingGroup.POST("/", bookingHandler.BookTour)
	bookingGroup.DELETE("/:id", bookingHandler.CancelBooking)

	paymentGroup := api.Group("/payments")
	paymentGroup.POST("/vnpay", middlewares.AuthMiddleware(), vnpayHandler.CreatePaymentUrl)
	paymentGroup.GET("/vnpay/callback", vnpayHandler.VnpayReturn)

	return r
}
