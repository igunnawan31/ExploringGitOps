package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/igunnawan31/enterprise-user/internal/config"
	"github.com/igunnawan31/enterprise-user/internal/handler"
	"github.com/igunnawan31/enterprise-user/internal/repository"
	"github.com/igunnawan31/enterprise-user/internal/service"
	"github.com/igunnawan31/enterprise-user/internal/validator"
)

func main() {
	db := config.ConnectDatabase()

	// Dependency Injection
	// give "db" to repository
	userRepository := repository.NewUserRepository(db)

	// give "userRepository" to service
	userService := service.NewUserService(userRepository)

	// give "userService" to handler
	userHandler := handler.NewUserHandler(userService)

	validator.InitValidators()

	// Set up Gin Framework
	router := gin.Default()

	// Define Routes
	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api)

	// Untuk debugging, untuk memastikan server berjalan
	log.Println("Server is running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server : %v", err)
	}

	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "ping",
	// 	})
	// })

	// router.Run()
}
