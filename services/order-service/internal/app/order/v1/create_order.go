package v1

import (
	"context"

	"order-service/internal/application/service/order/create_order"
	"order-service/internal/pkg/convert"
	orderV1 "order-service/internal/pkg/pb/order-service/order/v1"

	"github.com/gofrs/uuid"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error) {
	dto := create_order.NewCreateRequest(
		uuid.FromStringOrNil(req.GetOrder().GetAdId()),
		req.GetOrder().GetCategory(),
		uuid.FromStringOrNil(req.GetOrder().GetClientId()),
		convert.MoneyToDecimal(req.GetOrder().GetPrice()),
		// idempotency params
		req.GetIdempotencyKey(),
		req.GetOrdersVersion(),
	)

	details, err := i.services.CreateOrder.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		Details: &orderV1.CreateOrderResponse_Details{
			OrderId:     details.OrderID.String(),
			OrderStatus: details.Status,
		},
		OrdersVersion: details.OrdersVersion,
	}, nil
}
