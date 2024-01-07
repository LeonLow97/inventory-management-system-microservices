package kafkago

import (
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	BrokerAddress string
	TopicName     string
}

func NewKafkaConfig(brokerAddress string, topicName string) *KafkaConfig {
	return &KafkaConfig{
		BrokerAddress: brokerAddress,
		TopicName:     topicName,
	}
}

type Segmentio struct {
	topicConfigs []kafka.TopicConfig // to store multiple topics
}

func NewSegmentio() *Segmentio {
	return &Segmentio{}
}

// AddTopicConfig adds a new topic info apache kafka
func (s *Segmentio) AddTopicConfig(topic string, partitions int, replicationFactor int) {
	// add a new topic configuration
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}
	s.topicConfigs = append(s.topicConfigs, topicConfig)
}

// CreateTopics creates multiple topics into the broker, assuming (auto.create.topics.enable='false')
// where topics not in the broker are not auto created
func (s *Segmentio) CreateTopics(broker string) (*kafka.Conn, *kafka.Conn, error) {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return nil, nil, err
	}

	controller, err := conn.Controller()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	// controllerConn is the kafka controller for managing administrative tasks like creating topics
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	// create multiple topics
	err = controllerConn.CreateTopics(s.topicConfigs...)
	if err != nil {
		conn.Close()
		controllerConn.Close()
		return nil, nil, err
	}

	return conn, controllerConn, nil
}
