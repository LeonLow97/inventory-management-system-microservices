package kafkago

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic     = "DECREMENT_INVENTORY"
	partition = 0
	broker    = "broker:9092"
)

func Consumer() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, partition)
	if err != nil {
		log.Fatalf("Failed to dial leader: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal("failed to close connection:", err)
		}
	}()

	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatalf("Failed to set read deadline: %v", err)
	}

	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	defer func() {
		if err := batch.Close(); err != nil {
			log.Fatal("failed to close batch:", err)
		}
	}()

	b := make([]byte, 10e3) // 10 KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

}
