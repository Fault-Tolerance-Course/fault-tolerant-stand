package order_events

import (
	"context"
	"encoding/json"

	event_router "payment-service/internal/pkg/event-router"
	"payment-service/internal/pkg/inbox/message"
)

type BaseEvent struct {
	EventType string          `json:"event_type"`
	EntityID  string          `json:"entity_id"`
	Payload   json.RawMessage `json:"payload"`
}

type Handler struct {
	topic     string
	batchSize int

	router *event_router.EventRouter[string, []byte]
}

func NewHandler(topic string, batchSize int, router *event_router.EventRouter[string, []byte]) *Handler {
	return &Handler{
		topic:     topic,
		batchSize: batchSize,
		router:    router,
	}
}

func (h *Handler) Handle(ctx context.Context, messages message.Messages) error {
	for _, msg := range messages {
		var base BaseEvent
		err := json.Unmarshal(msg.GetBody(), &base)
		if err != nil {
			return err
		}

		err = h.router.HandleEvent(ctx, base.EventType, base.Payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) Topic() string {
	return h.topic
}

func (h *Handler) BatchSize(_ context.Context) int {
	return h.batchSize
}
