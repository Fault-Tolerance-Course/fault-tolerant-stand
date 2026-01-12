package order

import (
	"context"

	orderV1 "api-gateway/internal/pkg/pb/api-gateway/order/v1"
	externalV1 "api-gateway/internal/pkg/pb/external/order-service/order/v1"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error) {
	response, err := i.external.CreateOrder(ctx, &externalV1.CreateOrderRequest{
		Order: &externalV1.CreateOrderRequest_Order{
			AdId:     req.GetOrder().GetAdId(),
			Category: req.GetOrder().GetCategory(),
			ClientId: req.GetOrder().GetClientId(),
			Price:    req.GetOrder().GetPrice(),
		},
		IdempotencyKey: req.GetIdempotencyKey(),
		OrdersVersion:  req.GetOrdersVersion(),
	})
	if err != nil {
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		Details: &orderV1.CreateOrderResponse_Details{
			OrderId:     response.GetDetails().GetOrderId(),
			OrderStatus: response.GetDetails().GetOrderStatus(),
		},
		OrdersVersion: response.GetOrdersVersion(),
	}, nil
}
