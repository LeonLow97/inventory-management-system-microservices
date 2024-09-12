package kafkago

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func (s *Segmentio) Producer(brokerAddress, topic string, messages []kafka.Message) error {
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

	// produce messages to the topic
	if err := w.WriteMessages(context.Background(), messages...); err != nil {
		log.Println("failed to write messages:", err)
		return err
	}

	return nil
}
