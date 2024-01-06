package main

import (
	"log"

	kafkago "github.com/LeonLow97/internal/kafka"
)

const (
	topicDecrementInventory = "DECREMENT_INVENTORY"
	brokerAddress           = "broker:9092"
)

func main() {
	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

	segmentioInstance.AddTopicConfig(topicDecrementInventory, 1, 1)
	conn, controllerConn, err := segmentioInstance.CreateTopics(brokerAddress)
	if err != nil {
		log.Fatalln("Unable to create kafka topics", err)
	}
	log.Println("Successfully created kafka topics!")
	defer conn.Close()
	defer controllerConn.Close()

	go func() {
		if err := segmentioInstance.Producer(brokerAddress, topicDecrementInventory); err != nil {
			log.Printf("failed to produce message for %s topic: %v\n", topicDecrementInventory, err)
		}
	}()

	select {}
}
