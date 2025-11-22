package messagebus

import (
	"order-service/internal/infrastructure/messagebus/producer"
	"order-service/internal/pkg/closer"
)

type Registry struct {
	Producers producer.Producers
}

func NewRegistry() *Registry {
	registry := &Registry{
		Producers: producer.NewProducers(),
	}

	closer.Add(registry.Producers.OrderEvents.Close)
	return registry
}
