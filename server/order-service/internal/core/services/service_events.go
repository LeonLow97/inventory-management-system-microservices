package services

import (
	"github.com/LeonLow97/internal/ports"
)

type ServiceEvents interface {
	ConsumeUpdateInventoryEvent(brokerAddress, consumeTopic string)
}

type serviceEvents struct {
	repo ports.Repository
}

func NewServiceEvents(repo ports.Repository) ServiceEvents {
	return &serviceEvents{
		repo: repo,
	}
}

func (s *serviceEvents) ConsumeUpdateInventoryEvent(brokerAddress, topic string) {
	// TODO: Take from config .yaml
	// const (
	// 	topicUpdateOrderStatus = "update-order-status"
	// 	brokerAddress          = "broker:9092"
	// )

	s.repo.ConsumeOrderStatus(brokerAddress, topic)
}
