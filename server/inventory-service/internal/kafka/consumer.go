package kafkago

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// Consumer reads messages from a Kafka topic
func (s *Segmentio) Consumer(brokerAddress, topic string) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})
	defer r.Close()

	// read messages indefinitely
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("error reading message in consumer for topic", topic)
			break
		}
		fmt.Printf("Received message: %s\n", string(m.Value))
	}

	return nil
}
