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
	s.repo.ConsumeOrderStatus(brokerAddress, topic)
}
