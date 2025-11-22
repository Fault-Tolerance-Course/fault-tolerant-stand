package order

import (
	"context"

	commoddto "order-service/internal/application/dto"

	"order-service/internal/events/order_created"
	"order-service/internal/pkg/event"
)

type StorageInterface interface {
	CreateOrder(ctx context.Context, order *commoddto.Order) error
}

type Service struct {
	storage StorageInterface
	flusher event.Flusher
}

func NewService(storage StorageInterface, flusher event.Flusher) *Service {
	return &Service{storage: storage, flusher: flusher}
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (Details, error) {
	// создаем буфер событий
	buf, ctx := event.WithContext(ctx, s.flusher)

	order := commoddto.NewOrder(request.clientID, request.price, request.adID, request.category)

	// событие о создании заказа
	event.Add(ctx, order_created.New(order))

	// отправляем событие о создании
	if err := buf.Flush(ctx); err != nil {
		return Details{}, err
	}
	// создаем заказ
	if err := s.storage.CreateOrder(ctx, order); err != nil {
		return Details{}, err
	}

	return Details{
		OrderID: order.ID(),
		Status:  string(order.Status()),
	}, nil
}
