package event

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

const (
	eventTypeOrderCreated = "order-created"
)

type OrderCreated struct {
	ID        string    `json:"order_id"`
	Status    string    `json:"status"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	ClientID  string    `json:"client_id"`
	AdID      string    `json:"ad_id"`
}

// Deserialize переводит сообщение из одного представления в другое и валидирует
func Deserialize(message *sarama.ConsumerMessage) (OrderCreated, error) {
	var orderCreated OrderCreated
	err := json.Unmarshal(message.Value, &orderCreated)
	if err != nil {
		return OrderCreated{}, err
	}

	return orderCreated, nil
}
