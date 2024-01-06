package kafkago

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func (s *Segmentio) Producer(brokerAddress, topic string) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer func() {
		if err := w.Close(); err != nil {
			log.Fatalln("failed to close writer:", err)
		}
	}()

	if err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Test 1 Message"),
			Value: []byte("Decrease 1 count"),
		},
		kafka.Message{
			Key:   []byte("Test 2 Message"),
			Value: []byte("Decrease 100 count"),
		},
	); err != nil {
		log.Fatalln("failed to write messages:", err)
	}

	return nil
}
