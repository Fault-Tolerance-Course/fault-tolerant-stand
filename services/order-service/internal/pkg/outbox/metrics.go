package outbox

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "outbox"

var (
	// Векторные метрики, деление метрик происходит по топику
	canNotHandleMessageByOutbox *prometheus.CounterVec // Метрика количества неудачно обработанных с помощью outbox сообщений
	outboxMessagesWritten       *prometheus.CounterVec // Метрика количества записанных в outbox сообщений
	outboxMessagesProcessed     *prometheus.CounterVec // Метрика количества обработанных с помощью outbox сообщений

	once sync.Once
)

func registerMetrics() {
	once.Do(func() {
		canNotHandleMessageByOutbox = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "course_can_not_handler_message_by_outbox",
		}, []string{"topic"})

		outboxMessagesWritten = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "course_outbox_messages_written",
		}, []string{"topic"})
		outboxMessagesProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "course_outbox_messages_processed",
		}, []string{"topic"})

		prometheus.MustRegister(
			canNotHandleMessageByOutbox,
			outboxMessagesWritten,
			outboxMessagesProcessed,
		)
	})
}

func incCanNotHandleMessageByOutbox(topic string) {
	if canNotHandleMessageByOutbox != nil {
		canNotHandleMessageByOutbox.WithLabelValues(topic).Inc()
	}
}

func incOutboxMessagesWritten(topic string) {
	if outboxMessagesWritten != nil {
		outboxMessagesWritten.WithLabelValues(topic).Inc()
	}
}

func incOutboxMessagesProcessed(topic string, count float64) {
	if outboxMessagesProcessed != nil {
		outboxMessagesProcessed.WithLabelValues(topic).Add(count)
	}

}
