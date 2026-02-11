package order_created

import (
	"context"

	"payment-service/internal/application/service/payment"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const eventTypeOrderCreated = "order-created"

type PaymentService interface {
	Create(ctx context.Context, request payment.CreateRequest) error
}

type Handler struct {
	service PaymentService
}

func NewHandler(service PaymentService) *Handler {
	return &Handler{service: service}
}

// HandleEvent обработчик события
func (h *Handler) HandleEvent(ctx context.Context, payload []byte) error {
	event, err := deserialize(payload)
	if err != nil {
		return err
	}

	amount, err := decimal.NewFromString(event.Price)
	if err != nil {
		return err
	}

	return h.service.Create(ctx, payment.NewCreateRequest(
		uuid.MustParse(event.ID),
		event.ClientID,
		amount,
	))
}

// EventType возвращает тип ивента
func (h *Handler) EventType() string {
	return eventTypeOrderCreated
}
