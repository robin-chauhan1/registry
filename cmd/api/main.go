package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	AuthRoutes "ondc-registry/internals/api/routes/auth"
	RegistryRoutes "ondc-registry/internals/api/routes/registry"
	Userroutes "ondc-registry/internals/api/routes/user"
	"ondc-registry/internals/config/container"
	"ondc-registry/internals/config/database"
	"ondc-registry/internals/models"
	"ondc-registry/internals/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database configuration
	dbConfig := database.PostgresConfig{
		Host:            utils.GetEnv("DB_HOST", "localhost"),
		Port:            utils.GetEnv("DB_PORT", "5432"),
		User:            utils.GetEnv("DB_USER", "postgres"),
		Password:        utils.GetEnv("DB_PASSWORD", "postgres"),
		DBName:          utils.GetEnv("DB_NAME", "postgres"),
		SSLMode:         utils.GetEnv("DB_SSL_MODE", "disable"),
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
	}

	// Connect to database
	if err := dbConfig.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get database connection
	db := database.GetDB()
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}

	// Initialize dependency container
	container := container.NewContainer(db)

	// Ensure database connection is closed on application shutdown
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Run database migrations
	if err := dbConfig.Migrate(&models.User{}, models.Registry{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// API routes
	api := r.Group("/api")
	Userroutes.RegisterUserRoutes(api)
	AuthRoutes.AuthRoutes(api)
	RegistryRoutes.RegistryRoutes(api, container)

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
