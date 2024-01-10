package inventory

import (
	"encoding/json"
	"log"
	"time"

	"github.com/LeonLow97/pkg/kafkago"
	"github.com/segmentio/kafka-go"
)

func (s service) ConsumeKafkaUpdateInventoryCount() error {
	const (
		topicDecrementInventory = "update-inventory-count"
		topicUpdateOrderStatus  = "update-order-status"
		brokerAddress           = "broker:9092"
	)

	var (
		messageChan = make(chan interface{})
		errorChan   = make(chan error)
	)

	go s.segmentioInstance.Consumer(brokerAddress, topicDecrementInventory, messageChan, errorChan)
	for {
		select {
		case msg := <-messageChan:
			var orderEvent kafkago.OrderEvent
			if err := json.Unmarshal(msg.([]byte), &orderEvent); err != nil {
				log.Println("error unmarshaling message:", err)
				continue
			}

			// simulate order processing, in reality we may have tons of orders and the queue will take a long time
			time.Sleep(time.Second * 1)

			// TODO: run in database transaction
			// get inventory count and determine if sufficient for order
			dto := GetProductByIdDTO{
				UserID:    orderEvent.UserID,
				ProductID: orderEvent.ProductID,
			}
			product, err := s.repo.GetProductByID(dto)
			if err != nil {
				// stop processing messages
				log.Println("error getting product by id in kafka consumer", err)
				break
			}

			orderEventResp := kafkago.OrderEventResp{
				OrderUUID: orderEvent.OrderUUID,
				UserID:    orderEvent.UserID,
				ProductID: orderEvent.ProductID,
			}

			// insufficient inventory
			if product.Quantity < orderEvent.Quantity {
				// produce message with order uuid back to order microservice with status "FAILED"
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
				updateProduct := UpdateProductDTO{
					UserID:    orderEvent.UserID,
					ProductID: orderEvent.ProductID,
					Quantity:  finalQuantity,
				}
				if err := s.repo.UpdateProductByID(updateProduct); err != nil {
					log.Println("error updating product quantity when consuming message:", err)
					break
				}

				// produce message with order uuid back to order microservice with status "COMPLETED"
				orderEventResp.Status = "COMPLETED"
				orderEventResp.StatusReason = "Order completed!"
				orderEventResp.RemainingQuantity = finalQuantity
			}

			jsonData, err := json.Marshal(orderEvent)
			if err != nil {
				log.Println("error marshaling order event", err)
				break
			}
			updateOrderEvent := []kafka.Message{
				{
					Key:   []byte(orderEvent.OrderUUID),
					Value: []byte(jsonData),
				},
			}

			go func() {
				if err := s.segmentioInstance.Producer(brokerAddress, topicUpdateOrderStatus, updateOrderEvent); err != nil {
					log.Printf("failed to produce message for %s topic, order_uuid: %s, error: %v\n", topicDecrementInventory, orderEvent.OrderUUID, err)
				}
			}()
		case err := <-errorChan:
			// TODO: figure out what to do
			log.Println("error reading order", err)
		}
	}
}
