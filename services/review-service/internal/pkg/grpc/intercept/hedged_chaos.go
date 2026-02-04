package intercept

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

func HedgedDemoInterceptor() grpc.UnaryServerInterceptor {
	// методы, которые мы будем “мучать”
	methodMap := map[string]*uint64{
		"/review_service.review.v1.ReviewService/GetAdReview": new(uint64),
	}

	const delay = 1 * time.Second // задержка для "мучаемого" вызова
	const everyN = 5              // каждый N-й вызов задерживается

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		counter, ok := methodMap[info.FullMethod]
		if !ok {
			// метод не в списке -> просто вызываем хендлер
			return handler(ctx, req)
		}

		// увеличиваем счётчик атомарно
		reqNum := atomic.AddUint64(counter, 1)

		if int(reqNum)%everyN == 0 {
			slog.Info("[server] artificial delay applied",
				"method", info.FullMethod,
				"delay", delay,
			)

			select {
			case <-ctx.Done():
				slog.Error(fmt.Sprintf("request ctx err: %s", ctx.Err()))
				return nil, ctx.Err()
			case <-time.After(delay):
				return handler(ctx, req)
			}
		}

		slog.Info("[server] no delay", "method", info.FullMethod)

		return handler(ctx, req)
	}
}
