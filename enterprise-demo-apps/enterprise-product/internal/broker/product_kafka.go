package broker

import (
	"log"
	"os"

	"github.com/igunnawan31/enterprise-product/internal/service"
	"github.com/segmentio/kafka-go"
)

func StartKafkaProduct(productService service.ProductServiceInterface) error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")

	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}

	topic := "transactions-event"
	groupID := "enterprise-product-group"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e4,
		MaxBytes: 10e6,
	})

	log.Printf("Kafka Product started, listening to topic `%s`...", topic)

	go func() {
		defer reader.Close()

		for {
			// Membaca pesan masuk dari kafka dengan dua output: message dan error message
			// msg, err := reader.ReadMessage(context.Background())
			// if err != nil {
			// 	log.Printf("Error reading message from Kafka: %v", err)
			// 	break
			// }
		}
	}()

	return nil
}
