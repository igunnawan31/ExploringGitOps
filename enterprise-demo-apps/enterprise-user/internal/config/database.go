package config

import (
	"fmt"
	"log"
	"os"

	"github.com/igunnawan31/enterprise-user/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func ConnectDatabase() *gorm.DB {
// 	dsn := "host=db user=postgres password=rahasia dbname=db_users port=5432 sslmode=disable TimeZone=Asia/Jakarta"

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}

// 	err = db.AutoMigrate(&repository.User{})
// 	if err != nil {
// 		log.Fatalf("Failed to migrate database: %v", err)
// 	}

// 	log.Println("Database connection established")
// 	return db
// }

func ConnectDatabase() *gorm.DB {
	// 1. Membaca konfigurasi dari Environment Variables dengan nilai cadangan (fallback)
	host := getEnv("DB_HOST", "db") // Default 'db' untuk Docker Compose, atau 'localhost'
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "rahasia")
	dbname := getEnv("DB_NAME", "db_users")
	port := getEnv("DB_PORT", "5432")

	// 2. Menyusun DSN (Data Source Name) secara dinamis
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port,
	)

	// 3. Membuka koneksi database menggunakan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&repository.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
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
