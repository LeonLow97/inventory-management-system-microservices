package kafkago

type OrderEvent struct {
	OrderUUID string `json:"order_uuid"`
	ProductID int    `json:"product_id"`
	UserID    int    `json:"user_id"`
	Quantity  int    `json:"quantity"`
}

type InvEventConsume struct {
	OrderUUID         string `json:"order_uuid"`
	ProductID         int    `json:"product_id"`
	UserID            int    `json:"user_id"`
	Status            string `json:"status"`
	StatusReason      string `json:"status_reason"`
	RemainingQuantity int    `json:"remaining_quantity"`
}
