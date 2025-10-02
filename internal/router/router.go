package router

import (
	"project-simple/internal/config"
	"project-simple/internal/handler"
	"project-simple/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(cfg *config.Config, carHandler *handler.CarHandler, healthHandler *handler.HealthHandler) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router without default middleware
	router := gin.New()

	// Apply core middlewares (order matters!)
	router.Use(middleware.Recovery())           // Recover from panics
	router.Use(middleware.RequestID())          // Add request ID for tracing
	router.Use(middleware.Logger())             // Log requests
	router.Use(middleware.SecurityHeaders())    // Add security headers
	router.Use(middleware.CORS(cfg.Server.AllowedOrigins)) // CORS with configured origins
	router.Use(middleware.RequestSizeLimit(1 << 20))       // Limit request body to 1MB
	router.Use(middleware.ErrorHandler())       // Handle errors

	// Apply rate limiting (100 requests per minute per IP)
	router.Use(middleware.RateLimit(100, time.Minute))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check (no rate limit needed)
		v1.GET("/health", healthHandler.HealthCheck)

		// Car routes
		cars := v1.Group("/cars")
		{
			cars.POST("", carHandler.CreateCar)
			cars.GET("", carHandler.GetAllCars)
			cars.GET("/:id", carHandler.GetCarByID)
			cars.PUT("/:id", carHandler.UpdateCar)
			cars.DELETE("/:id", carHandler.DeleteCar)
		}
	}

	return router
}
