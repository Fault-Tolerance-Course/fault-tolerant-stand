package ratelimit

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/samber/lo"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// setHeader заголовки проставляются в тот момент, когда начинает отправляться
// тело ответа
func setHeader(ctx context.Context, l *rate.Limiter) {
	md := metadata.Pairs(
		"x-rate-limit", fmt.Sprintf("%d", int(l.Limit())),
		"x-rate-limit-burst", fmt.Sprintf("%d", l.Burst()),
		"x-rate-limit-remaining", fmt.Sprintf("%d", int(math.Max(0, l.Tokens()))),
	)
	_ = grpc.SetHeader(ctx, md)
}

// setTrailers трейлеры проставляются уже по результатам rpc
func setTrailers(ctx context.Context, delay time.Duration) {
	result := lo.Ternary(delay != 0, delay, -1)
	_ = grpc.SetTrailer(ctx, metadata.Pairs(
		"x-retry-after", fmt.Sprintf("%s", result.String()),
		"x-rate-limit-reset", fmt.Sprintf("%s", result.String()),
	))
}
