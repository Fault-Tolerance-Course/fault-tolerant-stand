package create_order

import (
	"context"

	"order-service/internal/domain/entity"
	"order-service/internal/domain/value_object"
	"order-service/internal/events/order_created"
	"order-service/internal/pkg/ordererror"
	"order-service/internal/pkg/transaction"

	"order-service/internal/pkg/event"

	"github.com/gofrs/uuid"
)

type Storage interface {
	CreateOrder(ctx context.Context, order *entity.Order, key string) error

	GetOrdersVersion(ctx context.Context, clientID uuid.UUID) (value_object.Version, error)

	UpdateOrdersVersion(ctx context.Context, clientID uuid.UUID, version value_object.Version) error
}

type Service struct {
	storage Storage
	flusher event.Flusher
}

func NewService(storage Storage, flusher event.Flusher) *Service {
	return &Service{storage: storage, flusher: flusher}
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (Details, error) {
	// создаем буфер событий
	buf, ctx := event.WithContext(ctx, s.flusher)

	order := entity.NewOrder(request.clientID, request.price, request.adID, request.category)

	// событие о создании заказа
	event.Add(ctx, order_created.New(order))

	// получаем текущую версию списка заказов пользователя
	version, err := s.storage.GetOrdersVersion(ctx, request.clientID)
	if err != nil {
		return Details{}, err
	}

	// если версии списка заказов отличаются - не даем делать заказ!
	if request.ordersVersion != version.Current() {
		return Details{}, ordererror.NewConflictErr("unable to create order due to version conflict")
	}

	// открываем транзакцию
	err = transaction.Exec(ctx, func(ctx context.Context) error {
		// пробуем создать заказ
		err = s.storage.CreateOrder(ctx, order, request.idempotencyKey)
		if err != nil {
			return err
		}

		// если заказ был создан и проблем не наблюдалось - повышаем версию
		version.Inc()

		return s.storage.UpdateOrdersVersion(ctx, request.clientID, version)
	})
	if err != nil {
		return Details{}, err
	}

	// отправляем событие о создании
	if err = buf.Flush(ctx); err != nil {
		return Details{}, err
	}

	return Details{
		OrderID:       order.ID(),
		Status:        string(order.Status()),
		OrdersVersion: version.Current(),
	}, nil
}
