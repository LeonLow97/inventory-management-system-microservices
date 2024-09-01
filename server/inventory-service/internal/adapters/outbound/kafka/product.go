package outbound_kafka

import (
	"github.com/LeonLow97/internal/ports"
	"github.com/LeonLow97/pkg/kafkago"
	"github.com/segmentio/kafka-go"
)

type KafkaAdapter struct {
	segmentioInstance *kafkago.Segmentio
}

func NewKafkaAdapter(segmentioInstance *kafkago.Segmentio) ports.EventBus {
	return &KafkaAdapter{
		segmentioInstance: segmentioInstance,
	}
}

func (ka *KafkaAdapter) ConsumeOrderMessage(brokerAddress, topic string, messageChan chan interface{}, errorChan chan error) {
	go ka.segmentioInstance.Consumer(brokerAddress, topic, messageChan, errorChan)
}

func (ka *KafkaAdapter) ProduceOrderMessage(brokerAddress, topic string, message []kafka.Message) error {
	return ka.segmentioInstance.Producer(brokerAddress, topic, message)
}
