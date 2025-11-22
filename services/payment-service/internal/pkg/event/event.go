package event

import (
	"context"

	"payment-service/internal/pkg/pipe"
)

// Raw бинарное представление содержимого
type Raw []byte

// Events список событий
type Events []Event

// Event событие системы
type Event struct {
	// Ключ, тело и заголовки сообщения
	Key, Body, Headers Raw
	// Schema идентификатор схемы (топика) сообщения
	Schema string
}

// Add добавляет в буфер из контекста событие.
// Если в контексте не было буфера - panic
func Add(ctx context.Context, eventCallback pipe.Func[Events]) {
	Extract(ctx).Add(eventCallback)
}
