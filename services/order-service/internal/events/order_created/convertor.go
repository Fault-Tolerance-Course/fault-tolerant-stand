package order_created

import (
	"context"
	"encoding/json"
	"time"

	"order-service/config"
	"order-service/internal/domain/entity"

	"order-service/internal/pkg/event"

	"order-service/internal/pkg/pipe"
)

type OrderCreatedEvent struct {
	ID        string    `json:"order_id"`
	Status    string    `json:"status"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	ClientID  string    `json:"client_id"`
	AdID      string    `json:"ad_id"`
}

func New(order *entity.Order) pipe.Func[event.Events] {
	return func(ctx context.Context, events event.Events) (event.Events, error) {
		// формируем событие
		body, err := json.Marshal(&OrderCreatedEvent{
			ID:        order.ID().String(),
			Status:    string(order.Status()),
			Price:     order.Price().String(),
			ClientID:  order.ClientID().String(),
			AdID:      order.AdID().String(),
			CreatedAt: order.CreatedAt(),
		})
		if err != nil {
			return nil, err
		}

		// загорловки в сообщении любые
		headers := map[string]string{
			"x-app-name":   "order-service",
			"x-event-type": "order-created",
		}

		// это выносится в общие функции
		headersRaw, err := json.Marshal(headers)
		if err != nil {
			return nil, err
		}

		return append(events, event.Event{
			Key:     event.Raw(order.ID().String()),
			Body:    body,
			Headers: headersRaw,
			Schema:  config.OrderEventsTopic,
		}), nil
	}
}
