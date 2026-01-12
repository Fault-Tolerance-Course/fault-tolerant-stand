package v1

import (
	"context"

	"order-service/internal/domain/entity"
	"order-service/internal/pkg/convert"
	orderV1 "order-service/internal/pkg/pb/order-service/order/v1"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"
)

func (i *Implementation) GetActiveOrders(ctx context.Context, req *orderV1.GetActiveOrdersRequest) (*orderV1.GetActiveOrdersResponse, error) {
	oderList, err := i.services.GetActiveOrders.GetActive(ctx, uuid.FromStringOrNil(req.GetClientId()))
	if err != nil {
		return nil, err
	}

	return &orderV1.GetActiveOrdersResponse{
		OrdersVersion: oderList.Version.Current(),
		Orders: lo.Map(oderList.Orders, func(order *entity.Order, _ int) *orderV1.GetActiveOrdersResponse_Order {
			return &orderV1.GetActiveOrdersResponse_Order{
				OrderId:     order.ID().String(),
				OrderStatus: string(order.Status()),
				AdId:        order.AdID().String(),
				Category:    order.Category(),
				Price:       convert.DecimalToMoney(order.Price()),
			}
		}),
	}, nil
}
