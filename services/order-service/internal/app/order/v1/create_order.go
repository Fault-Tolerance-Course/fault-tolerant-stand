package v1

import (
	"context"

	"order-service/internal/application/service/order"
	"order-service/internal/pkg/convert"
	orderV1 "order-service/internal/pkg/pb/order-service/order/v1"

	"github.com/google/uuid"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error) {
	dto := order.NewCreateRequest(
		uuid.MustParse(req.GetOrder().GetAdId()),
		req.GetOrder().GetCategory(),
		req.GetOrder().GetClientId(),
		convert.MoneyToDecimal(req.GetOrder().GetPrice()),
	)

	details, err := i.services.Order.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		Details: &orderV1.CreateOrderResponse_Details{
			OrderId:     details.OrderID.String(),
			OrderStatus: details.Status,
		},
	}, nil
}
