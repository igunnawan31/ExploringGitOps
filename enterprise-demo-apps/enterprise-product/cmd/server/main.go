package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/igunnawan31/enterprise-product/internal/broker"
	"github.com/igunnawan31/enterprise-product/internal/config"
	"github.com/igunnawan31/enterprise-product/internal/handler"
	"github.com/igunnawan31/enterprise-product/internal/repository"
	"github.com/igunnawan31/enterprise-product/internal/service"
)

func main() {
	db := config.ConnectDatabase()

	// Dependency Injection
	// give "db" to repository
	productRepository := repository.NewProductRepository(db)

	// give "productRepository" to service
	productService := service.NewProductService(productRepository)

	// give "productService" to handler
	productHandler := handler.NewProductHandler(productService)

	// Set up Validator
	// validator.InitValidators()

	// Set up Broker
	if err := broker.StartKafkaProduct(productService); err != nil {
		log.Fatalf("Failed to start Kafka product: %v", err)
	}

	// Set up Gin Framework
	router := gin.Default()

	// Define Routes
	api := router.Group("/api/v1")
	productHandler.RegisterRoutes(api)

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
