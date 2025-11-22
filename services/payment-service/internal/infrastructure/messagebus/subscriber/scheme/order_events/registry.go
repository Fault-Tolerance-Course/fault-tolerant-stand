package order_events

import (
	"payment-service/internal/application/service"
	"payment-service/internal/infrastructure/messagebus/subscriber/scheme/order_events/order_created"
)

type Registry struct {
	OrderCreated *order_created.MessageHandler
}

func NewRegistry(services *service.Registry) *Registry {
	return &Registry{
		OrderCreated: order_created.NewMessageHandler(services.Payment),
	}
}
