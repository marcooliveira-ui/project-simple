package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project-simple/internal/config"
	"project-simple/internal/handler"
	"project-simple/internal/infrastructure/database"
	"project-simple/internal/repository"
	"project-simple/internal/router"
	"project-simple/internal/service"
	"syscall"
	"time"

	_ "project-simple/docs" // Import swagger docs
)

// @title Car Management API
// @version 1.0
// @description A RESTful API for managing cars with CRUD operations and pagination
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@carapi.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	// Load configuration
	cfg := config.Load()
	log.Println("Configuration loaded successfully")

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	carRepo := repository.NewCarRepository(db.DB)

	// Initialize services
	carService := service.NewCarService(carRepo)

	// Initialize handlers
	carHandler := handler.NewCarHandler(carService)
	healthHandler := handler.NewHealthHandler(db.DB)

	// Setup router
	r := router.SetupRouter(cfg, carHandler, healthHandler)

	// Configure HTTP server with timeouts
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:           serverAddr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", serverAddr)
		log.Printf("Swagger documentation available at http://localhost%s/swagger/index.html", serverAddr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server gracefully...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Close database connections
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connections closed successfully")
	}

	log.Println("Server exited successfully")
}
