package v1

import (
	"context"
	"sync/atomic"
	"time"

	reviewV1 "review/internal/pkg/pb/review-service/review/v1"
)

func (i *Implementation) UnstableMethod(ctx context.Context, req *reviewV1.UnstableMethodRequest) (*reviewV1.UnstableMethodResponse, error) {
	// Определяем задержку на основе номера запроса или параметра
	requestNum := atomic.AddUint64(&i.unstableReqCounter, 1)

	var delay time.Duration
	switch {
	case requestNum <= 100:
		delay = 10 * time.Millisecond
	case requestNum <= 1000:
		delay = 1000 * time.Millisecond
	default:
		delay = 5000 * time.Millisecond
	}

	// Ждем либо завершения задержки, либо отмены контекста
	select {
	case <-time.After(delay):
		return &reviewV1.UnstableMethodResponse{}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
