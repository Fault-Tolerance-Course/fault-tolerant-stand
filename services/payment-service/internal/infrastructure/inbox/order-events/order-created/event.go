package order_created

import (
	"encoding/json"
	"time"
)

type OrderCreated struct {
	ID        string    `json:"order_id"`
	Status    string    `json:"status"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	ClientID  string    `json:"client_id"`
	AdID      string    `json:"ad_id"`
}

func deserialize(payload []byte) (OrderCreated, error) {
	var orderCreated OrderCreated
	err := json.Unmarshal(payload, &orderCreated)
	if err != nil {
		return OrderCreated{}, err
	}

	return orderCreated, nil
}
