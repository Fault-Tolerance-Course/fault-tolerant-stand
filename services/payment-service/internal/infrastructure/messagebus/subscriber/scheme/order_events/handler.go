package order_events

import (
	"context"
	"encoding/json"

	"payment-service/internal/pkg/inbox"
	"payment-service/internal/pkg/inbox/message"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type MessageHandler struct {
	inbox *inbox.Inbox
}

func NewMessageHandler(inbox *inbox.Inbox) *MessageHandler {
	return &MessageHandler{inbox: inbox}
}

func (h *MessageHandler) Handle(ctx context.Context, _ sarama.ConsumerGroupSession, kafkaMessage *sarama.ConsumerMessage) error {
	var base BaseEvent
	if err := json.Unmarshal(kafkaMessage.Value, &base); err != nil {
		return err
	}

	msg := message.NewMessage(
		base.EntityID,
		correlationID(kafkaMessage),
		kafkaMessage.Topic,
		kafkaMessage.Key,
		kafkaMessage.Value,
		kafkaMessage.Timestamp,
	)

	return h.inbox.SaveMessages(ctx, message.Messages{msg})
}

func correlationID(msg *sarama.ConsumerMessage) string {
	for _, h := range msg.Headers {
		if string(h.Key) == correlationIDHeader {
			return string(h.Value)
		}
	}
	return uuid.New().String()
}
