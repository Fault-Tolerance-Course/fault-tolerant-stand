package v1

import (
	orderV1 "order-service/internal/pkg/pb/order-service/order/v1"

	"order-service/internal/application/service"
)

type Implementation struct {
	orderV1.UnimplementedOrderServiceServer
	services *service.Registry

	requestCounter uint64 // atomic counter
}

func NewOrderService(services *service.Registry) *Implementation {
	return &Implementation{services: services}
}
