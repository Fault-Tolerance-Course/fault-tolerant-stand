package v1

import (
	"context"
	"sync/atomic"

	"order-service/internal/application/service/order/create_order"
	"order-service/internal/pkg/convert"
	orderV1 "order-service/internal/pkg/pb/order-service/order/v1"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error) {
	// увеличиваем атомарный счётчик
	count := atomic.AddUint64(&i.requestCounter, 1)
	switch {
	// каждый второй запрос
	case count%2 == 0:
		// если системы перегружена - отдает трейлеры с push back
		return trailers(ctx, true)
		// каждый третий запрос
	case count%3 == 0:
		// если системы перегружена - отдает трейлеры с push back
		return trailers(ctx, false)
	}

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

func trailers(ctx context.Context, noRetry bool) (*orderV1.CreateOrderResponse, error) {
	delay := lo.Ternary(noRetry, "-1", "1000")
	// каждый третий запрос → pushback 500ms
	md := metadata.Pairs("grpc-retry-pushback-ms", delay)
	grpc.SetTrailer(ctx, md)

	return nil, status.Error(codes.Unavailable, "system overloaded, retry later")
}
