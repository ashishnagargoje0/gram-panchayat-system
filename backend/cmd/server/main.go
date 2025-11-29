// cmd/server/main.go
package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gram-panchayat/internal/database"
	"gram-panchayat/internal/handlers"
	"gram-panchayat/internal/middleware"
	"gram-panchayat/internal/repository"
	"gram-panchayat/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db := database.InitDB()
	
	// Run migrations
	database.RunMigrations(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	complaintRepo := repository.NewComplaintRepository(db)
	propertyRepo := repository.NewPropertyRepository(db)
	noticeRepo := repository.NewNoticeRepository(db)
	meetingRepo := repository.NewMeetingRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	applicationService := service.NewApplicationService(applicationRepo)
	complaintService := service.NewComplaintService(complaintRepo)
	propertyService := service.NewPropertyService(propertyRepo)
	noticeService := service.NewNoticeService(noticeRepo)
	meetingService := service.NewMeetingService(meetingRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	complaintHandler := handlers.NewComplaintHandler(complaintService)
	propertyHandler := handlers.NewPropertyHandler(propertyService)
	noticeHandler := handlers.NewNoticeHandler(noticeService)
	meetingHandler := handlers.NewMeetingHandler(meetingService)
	dashboardHandler := handlers.NewDashboardHandler(userService, applicationService, complaintService)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Request logger
	r.Use(middleware.Logger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.PUT("/auth/profile", authHandler.UpdateProfile)
			protected.POST("/auth/change-password", authHandler.ChangePassword)

			// Dashboard
			protected.GET("/dashboard/admin", middleware.RoleMiddleware("admin"), dashboardHandler.GetAdminDashboard)
			protected.GET("/dashboard/citizen", dashboardHandler.GetCitizenDashboard)

			// Services/Applications
			applications := protected.Group("/services")
			{
				applications.POST("/apply", applicationHandler.CreateApplication)
				applications.GET("/applications", applicationHandler.GetUserApplications)
				applications.GET("/applications/:id", applicationHandler.GetApplication)
				applications.GET("/my-applications", applicationHandler.GetMyApplications)
			}

			// Complaints
			complaints := protected.Group("/complaints")
			{
				complaints.POST("", complaintHandler.CreateComplaint)
				complaints.GET("", complaintHandler.GetComplaints)
				complaints.GET("/:id", complaintHandler.GetComplaint)
				complaints.POST("/:id/comments", complaintHandler.AddComment)
			}

			// Property Tax
			properties := protected.Group("/property-tax")
			{
				properties.GET("/properties", propertyHandler.GetProperties)
				properties.POST("/properties", propertyHandler.CreateProperty)
				properties.GET("/:propertyId/bills", propertyHandler.GetBills)
				properties.POST("/:propertyId/payment", propertyHandler.MakePayment)
				properties.GET("/payment-history", propertyHandler.GetPaymentHistory)
			}

			// Notices (public read, admin write)
			notices := protected.Group("/notices")
			{
				notices.GET("", noticeHandler.GetNotices)
				notices.GET("/:id", noticeHandler.GetNotice)
				notices.POST("", middleware.RoleMiddleware("admin"), noticeHandler.CreateNotice)
				notices.PUT("/:id", middleware.RoleMiddleware("admin"), noticeHandler.UpdateNotice)
				notices.DELETE("/:id", middleware.RoleMiddleware("admin"), noticeHandler.DeleteNotice)
			}

			// Meetings
			meetings := protected.Group("/meetings")
			{
				meetings.GET("", meetingHandler.GetMeetings)
				meetings.GET("/:id", meetingHandler.GetMeeting)
				meetings.POST("", middleware.RoleMiddleware("admin"), meetingHandler.CreateMeeting)
				meetings.POST("/:id/minutes", middleware.RoleMiddleware("admin"), meetingHandler.AddMinutes)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("/users", userHandler.GetUsers)
				admin.GET("/users/:id", userHandler.GetUser)
				admin.PUT("/users/:id", userHandler.UpdateUser)
				admin.DELETE("/users/:id", userHandler.DeleteUser)
				
				admin.PUT("/applications/:id/status", applicationHandler.UpdateStatus)
				admin.PUT("/complaints/:id", complaintHandler.UpdateComplaint)
			}
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}