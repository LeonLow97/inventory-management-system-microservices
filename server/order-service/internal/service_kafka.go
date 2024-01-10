package order

import (
	"encoding/json"
	"log"

	"github.com/LeonLow97/pkg/kafkago"
)

func (s service) ConsumeKafkaUpdateOrderStatus() error {
	const (
		topicUpdateOrderStatus = "update-order-status"
		brokerAddress          = "broker:9092"
	)

	var (
		messageChan = make(chan interface{})
		errorChan   = make(chan error)
	)

	go s.segmentioInstance.Consumer(brokerAddress, topicUpdateOrderStatus, messageChan, errorChan)
	for {
		select {
		case msg := <-messageChan:
			var invEventConsume kafkago.InvEventConsume
			if err := json.Unmarshal(msg.([]byte), &invEventConsume); err != nil {
				log.Println("error unmarshaling message:", err)
				continue
			}

			req := UpdateOrderDTO{
				OrderUUID:    invEventConsume.OrderUUID,
				Status:       invEventConsume.Status,
				StatusReason: invEventConsume.StatusReason,
			}
			if err := s.repo.UpdateOrderByUUID(req); err != nil {
				log.Printf("error updating order for order uuid %s with error %v\n", req.OrderUUID, err)
			}
		case err := <-errorChan:
			log.Println("error reading inventory message", err)
		}
	}
}
