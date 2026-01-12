package circuit

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCircuitIsOpen = errors.New("circuit is open")

	ErrTooManyRequests = errors.New("too many requests through circuit")
)

// errorCodeMap отображение строк -> codes.Code
var errorCodeMap = map[string]codes.Code{
	"Unknown":           codes.Unknown,
	"DeadlineExceeded":  codes.DeadlineExceeded,
	"ResourceExhausted": codes.ResourceExhausted,
	"Internal":          codes.Internal,
	"Unavailable":       codes.Unavailable,
}

func triggerOnError(err error, codes []string, codesSet map[codes.Code]struct{}) bool {
	if err == nil {
		return false
	}

	// Пытаемся получить gRPC-статус
	st, ok := status.FromError(err)
	if !ok {
		// Не gRPC-ошибка — считаем за сбой, если список кодов не задан
		return len(codes) == 0
	}

	code := st.Code()

	// Если фильтрация не задана — считаем всё ошибками
	if len(codes) == 0 {
		return true
	}

	// Проверяем вхождение в список "плохих" кодов
	_, allowed := codesSet[code]
	return allowed
}
