package kafkago

type OrderEvent struct {
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	UserID    int `json:"user_id"`
	Quantity  int `json:"quantity"`
}

type InvEventConsume struct {
	OrderID           int    `json:"order_id"`
	ProductID         int    `json:"product_id"`
	UserID            int    `json:"user_id"`
	Status            string `json:"status"`
	StatusReason      string `json:"status_reason"`
	RemainingQuantity int    `json:"remaining_quantity"`
}
