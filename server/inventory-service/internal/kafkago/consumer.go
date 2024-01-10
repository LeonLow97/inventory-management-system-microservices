package kafkago

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

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
