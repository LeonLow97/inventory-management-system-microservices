package ports

import "github.com/segmentio/kafka-go"

// KafkaConsumer defines the interface for consuming Kafka messages
type EventBus interface {
	ConsumeOrderMessage(brokerAddress, topic string, messageChan chan interface{}, errorChan chan error)
	ProduceOrderMessage(brokerAddress, topic string, message []kafka.Message) error
}
