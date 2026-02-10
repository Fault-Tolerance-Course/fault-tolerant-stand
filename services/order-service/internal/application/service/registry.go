package service

import (
	"order-service/internal/application/service/order/create_order"
	"order-service/internal/application/service/order/get_active_orders"
	"order-service/internal/infrastructure/dal"
	"order-service/internal/pkg/outbox"
)

type Registry struct {
	CreateOrder     *create_order.Service
	GetActiveOrders *get_active_orders.Service
}

func NewRegistry(dal *dal.Registry, outbox *outbox.Outbox) *Registry {
	return &Registry{
		CreateOrder:     create_order.NewService(dal.OrderRepository, outbox),
		GetActiveOrders: get_active_orders.NewService(dal.OrderRepository),
	}
}
