package intercept

import (
	"context"
	"log/slog"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CircuitDemoInterceptor() grpc.UnaryServerInterceptor {
	// методы, которые мы будем “мучать”
	methodMap := map[string]*uint64{
		"/review_service.review.v1.ReviewService/GetAdReview": new(uint64),
	}

	const everyN = 3 // каждый N-й вызов будет ошибочным

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		counter, ok := methodMap[info.FullMethod]
		if !ok {
			// метод не в списке -> просто вызываем хендлер
			return handler(ctx, req)
		}

		// увеличиваем счётчик атомарно
		reqNum := atomic.AddUint64(counter, 1)

		// каждый N-й вызов -> генерируем ошибку
		if int(reqNum)%everyN == 0 {
			slog.Warn("[server] artificial failure injected",
				"method", info.FullMethod,
				"reqNum", reqNum,
			)
			return nil, status.Error(codes.Unavailable, "simulated failure for circuit breaker demo")
		}

		slog.Info("[server] normal execution", "method", info.FullMethod, "reqNum", reqNum)
		return handler(ctx, req)
	}
}
