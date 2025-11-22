package messagebus

import (
	"context"

	"payment-service/internal/application/service"
	"payment-service/internal/infrastructure/messagebus/subscriber"
	"payment-service/internal/infrastructure/messagebus/subscriber/scheme/order_events"
	"payment-service/internal/pkg/closer"
)

type handlers struct {
	OrderEvents *order_events.Registry
}

type Registry struct {
	handlers    handlers
	subscribers subscriber.Subscribers
}

func NewRegistry(services *service.Registry) *Registry {
	registry := &Registry{
		subscribers: subscriber.NewSubscribers(),
		handlers: handlers{
			OrderEvents: order_events.NewRegistry(services),
		},
	}

	closer.Add(registry.subscribers.OrderEvents.Close)
	return registry
}

func (r *Registry) Run(ctx context.Context) {
	// Можно добавить полноценный multiplexer для событий этого топика, но пока сосредоточимся на одном
	go r.subscribers.OrderEvents.Subscribe(ctx, r.handlers.OrderEvents.OrderCreated.Handle)
}
