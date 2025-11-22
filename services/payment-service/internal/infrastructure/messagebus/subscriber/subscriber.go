package subscriber

import (
	"payment-service/config"
	"payment-service/internal/pkg/msgbus/subscriber"
)

type Subscribers struct {
	OrderEvents *subscriber.MessageSubscriber
}

func NewSubscribers() Subscribers {
	return Subscribers{
		OrderEvents: subscriber.NewMessageSubscriber(config.OrderEventsTopic),
	}
}
