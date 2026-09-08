package broker

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type UserRegistrationEvent struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash,omitempty"`
	Provider     string `json:"provider"`
}

func PublishUserRegistrationEvent(eventData UserRegistrationEvent) error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	// Kondisi : Jika kakfkaBroker == "" maka Fallback nilai default Docker
	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}

	payload, err := json.Marshal(eventData)
	if err != nil {
		log.Printf("Failed to marshal kafka event: %v", err)
		return err
	}

	// Konfigurasi Writer Kafka untuk topik "user-registered"
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    "user-registered",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	// Mengirimkan pesan ke kafka dengan timeout 5 detik
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(eventData.Email),
		Value: payload,
	})

	if err != nil {
		log.Printf("Failed to publish message to Kafka: %v", err)
	} else {
		log.Printf("Successfully published message to Kafka for user: %s", eventData.Email)
	}

	return err
}
