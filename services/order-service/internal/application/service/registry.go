package service

import (
	"order-service/internal/application/service/order"
	"order-service/internal/infrastructure/dal"
	"order-service/internal/infrastructure/messagebus"
)

type Registry struct {
	Order *order.Service
}

func NewRegistry(dal *dal.Registry, messageBus *messagebus.Registry) *Registry {
	return &Registry{
		Order: order.NewService(dal.OrderRepository, messageBus.Producers.OrderEvents),
	}
}
