package order

import (
	orderV1 "api-gateway/internal/pkg/pb/api-gateway/order/v1"
	externalV1 "api-gateway/internal/pkg/pb/external/order-service/order/v1"
)

type Implementation struct {
	orderV1.UnimplementedOrderServiceServer

	external externalV1.OrderServiceClient
}

func NewOrderService(external externalV1.OrderServiceClient) *Implementation {
	return &Implementation{external: external}
}
