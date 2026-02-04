package v1

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	reviewV1 "review/internal/pkg/pb/review-service/review/v1"
)

func (i *Implementation) UnstableMethod(
	ctx context.Context,
	_ *reviewV1.UnstableMethodRequest,
) (*reviewV1.UnstableMethodResponse, error) {

	requestNum := atomic.AddUint64(&i.unstableReqCounter, 1)

	var (
		delay    time.Duration
		newStage uint32
	)

	switch {
	case requestNum <= 2000:
		delay = 10 * time.Millisecond
		newStage = stageFast

	case requestNum <= 5000 && requestNum > 2000:
		delay = 500 * time.Millisecond
		newStage = stageSlow500

	case requestNum <= 10000 && requestNum > 5000:
		delay = 1 * time.Second
		newStage = stageSlow1s

	case requestNum <= 15000 && requestNum > 10000:
		delay = 3 * time.Second
		newStage = stageSlow3s

	default:
		delay = 10 * time.Millisecond
		newStage = stageFastAgain
	}

	prev := atomic.LoadUint32(&i.stage)
	if prev != newStage {
		if atomic.CompareAndSwapUint32(&i.stage, prev, newStage) {
			slog.Info("UnstableMethod stage changed",
				"requestNum", requestNum,
				"delay", delay,
				"stage", newStage,
			)
		}
	}

	select {
	case <-time.After(delay):
		return &reviewV1.UnstableMethodResponse{}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
