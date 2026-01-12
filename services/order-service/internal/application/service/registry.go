package service

import (
	"order-service/internal/application/service/order/create_order"
	"order-service/internal/application/service/order/get_active_orders"
	"order-service/internal/infrastructure/dal"
	"order-service/internal/infrastructure/messagebus"
)

type Registry struct {
	CreateOrder     *create_order.Service
	GetActiveOrders *get_active_orders.Service
}

func NewRegistry(dal *dal.Registry, messageBus *messagebus.Registry) *Registry {
	return &Registry{
		CreateOrder:     create_order.NewService(dal.OrderRepository, messageBus.Producers.OrderEvents),
		GetActiveOrders: get_active_orders.NewService(dal.OrderRepository),
	}
}
