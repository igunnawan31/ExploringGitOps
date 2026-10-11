package config

import (
	"fmt"
	"log"
	"os"

	"github.com/igunnawan31/enterprise-product/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	host := getEnv("DB_HOST", "db") // Default 'db' untuk Docker Compose, atau 'localhost'
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "rahasia")
	dbname := getEnv("DB_NAME", "db_products")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&domain.Category{}, &domain.Product{})
	if err != nil {
		log.Fatalf("Failed to migrate database %v", err)
	}

	log.Println("Database connection successfully established via Environment Variables")
	return db
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
