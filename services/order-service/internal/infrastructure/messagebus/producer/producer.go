package producer

import (
	"order-service/internal/pkg/msgbus/producer"
)

type Producers struct {
	OrderEvents *producer.MessageProducer
}

func NewProducers() Producers {
	return Producers{
		OrderEvents: producer.NewMessageProducer(),
	}
}
