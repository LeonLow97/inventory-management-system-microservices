package services

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/LeonLow97/internal/pkg/kafkago"
	"github.com/LeonLow97/internal/ports"
	"github.com/segmentio/kafka-go"
)

type ServiceEvents interface {
	ConsumeUpdateInventoryEvent(brokerAddress, consumeTopic, produceTopic string) error
}

type serviceEvents struct {
	repo ports.Repository
}

func NewServiceEvents(repo ports.Repository) ServiceEvents {
	return &serviceEvents{
		repo: repo,
	}
}

func (s *serviceEvents) ConsumeUpdateInventoryEvent(brokerAddress, consumeTopic, produceTopic string) error {
	messageChan := make(chan interface{})
	errorChan := make(chan error)

	s.repo.ConsumeOrderMessage(brokerAddress, consumeTopic, messageChan, errorChan)
	for {
		select {
		case msg := <-messageChan:
			var orderEvent kafkago.OrderEvent
			if err := json.Unmarshal(msg.([]byte), &orderEvent); err != nil {
				log.Println("failed to unmarshal message when consuming event:", err)
				continue
			}

			product, err := s.repo.GetProductByID(orderEvent.UserID, orderEvent.ProductID)
			if err != nil {
				// stop processing the event
				log.Println("failed to get product by id when consuming event", err)
				break
			}

			orderEventResp := kafkago.OrderEventResp{
				OrderID:   orderEvent.OrderID,
				UserID:    orderEvent.UserID,
				ProductID: orderEvent.ProductID,
			}

			// insufficient inventory
			if product.Quantity < orderEvent.Quantity {
				// produce message with order id back to order microservice with status "FAILED"
				orderEventResp.Status = "FAILED"
				orderEventResp.RemainingQuantity = product.Quantity

				if product.Quantity == 0 {
					orderEventResp.StatusReason = "Sold Out"
				} else {
					orderEventResp.StatusReason = "Insufficient inventory count"
				}
			} else {
				// sufficient inventory
				// deduct inventory count and update quantity in database
				finalQuantity := product.Quantity - orderEvent.Quantity

				if err := s.repo.UpdateProductQuantityByID(finalQuantity, orderEvent.UserID, orderEvent.ProductID); err != nil {
					log.Println("error updating product quantity when consuming message:", err)
					break
				}

				// produce message with order id back to order microservice with status "COMPLETED"
				orderEventResp.Status = "COMPLETED"
				orderEventResp.StatusReason = "Order completed!"
				orderEventResp.RemainingQuantity = finalQuantity
			}

			jsonData, err := json.Marshal(orderEventResp)
			if err != nil {
				log.Println("error marshaling order event", err)
				break
			}
			updateOrderEvent := []kafka.Message{
				{
					Key:   []byte(strconv.Itoa(orderEvent.OrderID)),
					Value: []byte(jsonData),
				},
			}

			go func() {
				if err := s.repo.ProduceOrderMessage(brokerAddress, produceTopic, updateOrderEvent); err != nil {
					log.Printf("failed to produce message for %s topic, order_id: %d, error: %v\n", produceTopic, orderEvent.OrderID, err)
				}
			}()
		case err := <-errorChan:
			// TODO: figure out what to do
			log.Println("error reading order", err)
		}
	}
}
