package order

import (
	"context"

	orderV1 "api-gateway/internal/pkg/pb/api-gateway/order/v1"
	externalV1 "api-gateway/internal/pkg/pb/external/order-service/order/v1"

	"github.com/samber/lo"
)

func (i *Implementation) GetActiveOrders(ctx context.Context, req *orderV1.GetActiveOrdersRequest) (*orderV1.GetActiveOrdersResponse, error) {
	response, err := i.external.GetActiveOrders(ctx, &externalV1.GetActiveOrdersRequest{
		ClientId: req.GetClientId(),
	})
	if err != nil {
		return nil, err
	}

	return &orderV1.GetActiveOrdersResponse{
		Orders: lo.Map(response.GetOrders(), func(order *externalV1.GetActiveOrdersResponse_Order, _ int) *orderV1.GetActiveOrdersResponse_Order {
			return &orderV1.GetActiveOrdersResponse_Order{
				OrderId:     order.GetOrderId(),
				OrderStatus: order.GetOrderStatus(),
				AdId:        order.GetAdId(),
				Category:    order.GetCategory(),
				Price:       order.GetPrice(),
			}
		}),
		OrdersVersion: response.GetOrdersVersion(),
	}, nil
}
