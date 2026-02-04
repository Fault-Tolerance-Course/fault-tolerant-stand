package intercept

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func RetryDemoInterceptor() grpc.UnaryServerInterceptor {
	var requestCounter uint64

	// методы, которые мы будем “мучать”
	methodMap := map[string]struct{}{
		"/order_service.order.v1.OrderService/CreateOrder": {},
	}

	// список сценариев
	scenarioList := []string{"success", "push_back_reject", "push_back_delay", "retryable_fail"}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// если метод не выбран → просто хендлер
		if _, ok := methodMap[info.FullMethod]; !ok {
			return handler(ctx, req)
		}

		// читаем attempt
		attempt := 0
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get("x-retry-attempt"); len(values) > 0 {
				a, _ := strconv.Atoi(values[0])
				attempt = a
			}
		}

		// читаем max_attempts
		maxAttempts := 3
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get("x-retry-max-attempts"); len(values) > 0 {
				a, _ := strconv.Atoi(values[0])
				maxAttempts = a
			}
		}

		var scenario string

		if attempt == 0 {
			atomic.AddUint64(&requestCounter, 1)
		}

		reqNum := int(atomic.LoadUint64(&requestCounter))
		idx := (reqNum - 1) % len(scenarioList)
		scenario = scenarioList[idx]

		slog.Info(fmt.Sprintf("[server] method=%s reqNum=%d attempt=%d scenario=%s",
			info.FullMethod, atomic.LoadUint64(&requestCounter), attempt, scenario))

		switch scenario {
		case "success":
			// сразу пропускаем в хендлер → успех
			return handler(ctx, req)

		case "push_back_reject":
			if attempt == 0 {
				// сервер говорит клиенту НЕ ретраить
				grpc.SetTrailer(ctx, metadata.Pairs("grpc-retry-pushback-ms", "-1"))
				return nil, status.Error(codes.Unavailable, "pushback reject")
			}
			// если attempt > 0, просто хендлер (flow уже закончился)
			return handler(ctx, req)

		case "push_back_delay":
			if attempt == 0 {
				// первый attempt → пушбек с задержкой
				grpc.SetTrailer(ctx, metadata.Pairs("grpc-retry-pushback-ms", "2000"))
				return nil, status.Error(codes.Unavailable, "pushback delay")
			}
			// повторная попытка клиента → хендлер → успех
			return handler(ctx, req)

		case "retryable_fail":
			// отдаём ошибки до достижения maxRetryableAttempts
			if attempt < maxAttempts-1 {
				return nil, status.Error(codes.Unavailable, "retryable failure")
			}
			// последний attempt → финальная ошибка для клиента
			return nil, status.Error(codes.Unavailable, "final failure after retries")
		}

		// safety fallback
		return handler(ctx, req)
	}
}
