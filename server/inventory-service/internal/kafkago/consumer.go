package kafkago

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// TODO: Shift this orderevent to the service (business logic) area
type OrderEvent struct {
	Action    string `json:"action"`
	ProductID int    `json:"product_id"`
	UserID    int    `json:"user_id"`
	Quantity  int    `json:"quantity"`
}

func (s *Segmentio) Consumer(brokerAddress, topic string, messageChan chan interface{}, errorChan chan error) {
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
			errorChan <- err
			continue
		}

		// send raw Kafka message to messageChan
		messageChan <- m.Value
	}
}

// Consumer reads messages from a Kafka topic
func (s *Segmentio) Consumer2(brokerAddress, topic string, messageChan chan OrderEvent, errorChan chan error) {
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
			errorChan <- err
			continue
		}

		var orderEvent OrderEvent
		err = json.Unmarshal(m.Value, &orderEvent)
		if err != nil {
			log.Println("error unmarshaling message:", err)
			errorChan <- err
			continue
		}

		messageChan <- orderEvent
	}
}
