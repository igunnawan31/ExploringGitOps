package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/igunnawan31/enterprise-auth/internal/config"
	"github.com/igunnawan31/enterprise-auth/internal/handler"
	"github.com/igunnawan31/enterprise-auth/internal/repository"
	"github.com/igunnawan31/enterprise-auth/internal/service"
)

func main() {
	config.ConnectDatabase()

	// Define Repo
	authRepo := repository.NewAuthRepository()

	// Define Service
	authService := service.NewAuthService(authRepo)

	// Define Handler
	authHandler := handler.NewAuthHandler(authService)

	// Define Validator

	// Set up Gin Framework
	router := gin.Default()

	// Define Routes
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.GET("/google/login", authHandler.GoogleLogin)
		authRoutes.GET("/google/callback", authHandler.GoogleCallback)

		authRoutes.POST("/register", authHandler.RegisterManual)
		authRoutes.POST("/login", authHandler.LoginManual)
	}

	// Define Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Enterprise-Auth service is running smoothly on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
