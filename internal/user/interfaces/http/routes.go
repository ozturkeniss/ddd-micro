package http

import (
	"net/http"

	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/domain"
	"github.com/ddd-micro/internal/user/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all user-related routes
func SetupRoutes(router *gin.Engine, userService *application.UserServiceCQRS, metrics *monitoring.PrometheusMetrics, tracer *monitoring.JaegerTracer) {
	// Add monitoring middlewares
	router.Use(monitoring.PrometheusMiddleware(metrics))
	router.Use(monitoring.JaegerMiddleware(tracer))

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create handler
	userHandler := NewUserHandler(userService, metrics)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// ========== PUBLIC ROUTES ==========
		// No authentication required
		publicUsers := v1.Group("/users")
		{
			publicUsers.POST("/register", userHandler.Register)
			publicUsers.POST("/login", userHandler.Login)
			publicUsers.POST("/refresh-token", userHandler.RefreshToken)
		}

		// ========== USER ROUTES (User + Admin can access) ==========
		// Requires authentication as User or Admin
		userRoutes := v1.Group("/users")
		userRoutes.Use(AuthMiddleware(userService))
		userRoutes.Use(RequireRoles(domain.RoleUser, domain.RoleAdmin))
		{
			// Profile routes - authenticated users can manage their own profile
			userRoutes.GET("/profile", userHandler.GetProfile)
			userRoutes.PUT("/profile", userHandler.UpdateProfile)
			userRoutes.POST("/change-password", userHandler.ChangePassword)
		}

		// ========== ADMIN ROUTES (Only Admin can access) ==========
		// Requires authentication as Admin
		adminRoutes := v1.Group("/admin/users")
		adminRoutes.Use(AuthMiddleware(userService))
		adminRoutes.Use(RequireAdmin())
		{
			// User management - only admins
			adminRoutes.GET("", userHandler.ListUsers)
			adminRoutes.GET("/:id", userHandler.GetUserByID)
			adminRoutes.PUT("/:id", userHandler.UpdateUserByAdmin)
			adminRoutes.DELETE("/:id", userHandler.DeleteUser)
			adminRoutes.POST("/:id/assign-role", userHandler.AssignRole)
		}
	}
}
