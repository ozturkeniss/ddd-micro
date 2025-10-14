package http

import (
	"github.com/ddd-micro/internal/user/application"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all user-related routes
func SetupRoutes(router *gin.Engine, userService *application.UserService) {
	// Create handler
	userHandler := NewUserHandler(userService)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.POST("/refresh-token", userHandler.RefreshToken)
		}

		// Protected routes (require authentication)
		authenticated := v1.Group("/users")
		authenticated.Use(AuthMiddleware(userService))
		{
			// Profile routes
			authenticated.GET("/profile", userHandler.GetProfile)
			authenticated.PUT("/profile", userHandler.UpdateProfile)
			authenticated.POST("/change-password", userHandler.ChangePassword)

			// User management routes
			authenticated.GET("", userHandler.ListUsers)
			authenticated.GET("/:id", userHandler.GetUserByID)
			authenticated.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}

