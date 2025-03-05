package outbound

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/pkg/kafkago"
	"github.com/segmentio/kafka-go"
)

func (r *Repository) ProduceOrderMessage(brokerAddress, topic string, orderID int, productID, userID, orderQuantity int) error {
	orderEvent := kafkago.OrderEvent{
		OrderID:   orderID,
		ProductID: productID,
		UserID:    userID,
		Quantity:  orderQuantity,
	}

	jsonData, err := json.Marshal(orderEvent)
	if err != nil {
		return err
	}

	// produce an order to inventory microservice to update inventory count
	createOrderEvent := []kafka.Message{
		{
			Key:   []byte(strconv.Itoa(orderEvent.OrderID)),
			Value: []byte(jsonData),
		},
	}

	if err := r.segmentioInstance.Producer(brokerAddress, topic, createOrderEvent); err != nil {
		return err
	}
	return nil
}

func (r *Repository) ConsumeOrderStatus(brokerAddress, topic string) {
	var (
		messageChan = make(chan any)
		errorChan   = make(chan error)
	)

	go r.segmentioInstance.Consumer(brokerAddress, topic, messageChan, errorChan)
	for {
		select {
		case msg := <-messageChan:
			var invEvent kafkago.InvEventConsume
			if err := json.Unmarshal(msg.([]byte), &invEvent); err != nil {
				log.Printf("failed to unmarshal json when consuming events from inventory service with error: %v\n", err)
				continue
			}

			req := domain.Order{
				ID:           invEvent.OrderID,
				Status:       invEvent.Status,
				StatusReason: invEvent.StatusReason,
			}
			if err := r.UpdateOrderByID(req); err != nil {
				log.Printf("failed to update order id %d with error: %v\n", invEvent.OrderID, err)
			}
		case err := <-errorChan:
			log.Printf("failed to consume message from inventory service with error: %v\n", err)
		}
	}
}
