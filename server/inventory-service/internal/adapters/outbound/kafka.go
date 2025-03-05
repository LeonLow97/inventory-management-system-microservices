package outbound

import "github.com/segmentio/kafka-go"

func (r *repository) ConsumeOrderMessage(brokerAddress, topic string, messageChan chan any, errorChan chan error) {
	go r.segmentioInstance.Consumer(brokerAddress, topic, messageChan, errorChan)
}

func (r *repository) ProduceOrderMessage(brokerAddress, topic string, message []kafka.Message) error {
	return r.segmentioInstance.Producer(brokerAddress, topic, message)
}
