package order_created

import (
	"context"
	"fmt"
	"log/slog"

	"payment-service/internal/application/service/payment"
	"payment-service/internal/infrastructure/messagebus/subscriber/scheme/order_events/order_created/event"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentCreator interface {
	Create(ctx context.Context, request payment.CreateRequest) error
}

type MessageHandler struct {
	creator PaymentCreator
}

func NewMessageHandler(creator PaymentCreator) *MessageHandler {
	return &MessageHandler{creator: creator}
}

func (h *MessageHandler) Handle(ctx context.Context, _ sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) error {
	// десереализуем сообщение
	deserialized, err := event.Deserialize(message)
	if err != nil {
		slog.Error(fmt.Sprintf("Ошибка десереализации сообщения: %s", err.Error()))
		return err
	}

	amount, err := decimal.NewFromString(deserialized.Price)
	if err != nil {
		return err
	}

	return h.creator.Create(ctx, payment.NewCreateRequest(
		uuid.MustParse(deserialized.ID),
		deserialized.ClientID,
		amount,
	))
}
