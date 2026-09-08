package broker

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/igunnawan31/enterprise-user/internal/repository"
	"github.com/segmentio/kafka-go"
	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationEvent struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash,omitempty"`
	Provider     string `json:"provider"`
}

func StartKafkaUser(userRepo repository.UserRepositoryInterface) error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	// Kondisi : Jika kakfkaBroker == "" maka Fallback nilai default Docker
	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}

	topic := "user-registered"
	groupID := "enterprise-user-group"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e4,
		MaxBytes: 1063364,
	})

	log.Printf("Kafka User started, listening to topic `%s`...", topic)

	go func() {
		defer reader.Close()

		for {
			// Membaca pesan masuk dari kafka dengan dua output: message dan error message
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message from Kafka: %v", err)
				break
			}

			var event UserRegistrationEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Failed to unmarshal user registration event: %v", err)
				continue
			}

			var hashedPassword string
			if event.PasswordHash != "" {
				hashBytes, err := bcrypt.GenerateFromPassword([]byte(event.PasswordHash), 12)
				if err != nil {
					log.Printf("Failed to hash password from kafka event: %v", err)
					continue
				}
				hashedPassword = string(hashBytes)
			}

			// Mapping data event ke model database user
			newUser := repository.User{
				Name:         event.Name,
				Email:        event.Email,
				PasswordHash: hashedPassword,
			}

			if err := userRepo.Save(&newUser); err != nil {
				log.Printf("Failed to save registered user from kafka to database: %v", err)
			} else {
				log.Printf("Successfully saved user to database from Kafka event: %s", newUser.Email)
			}
		}
	}()

	return nil
}
