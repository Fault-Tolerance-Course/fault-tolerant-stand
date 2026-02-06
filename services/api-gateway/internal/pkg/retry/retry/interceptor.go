package retry

import (
	"context"
	"log/slog"
	"math"
	"math/rand/v2"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryClientInterceptor унарный клиентский перехватчик.
// Необходимо в опциях подключения установить grpc.WithDisableRetry(), чтобы использовать кастомный механизм
func (r *Retry) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		state, ok := r.getByMethod(method)
		if !ok || state == nil {
			slog.Debug("no retry mechanism found for method or disabled", "method", method)
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		var throttle = state.throttler

		var lastErr error
		for attempt := 0; attempt < state.cfg.MaxAttempts; attempt++ {
			// Проверяем, отменен ли контекст
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			attemptsTotal.WithLabelValues(method).Inc()

			// проверяем бюджет, если это retry
			if attempt != 0 {
				if throttle != nil && throttle.Throttle() {
					throttledTotal.WithLabelValues(method).Inc()
					return lastErr
				}
			}

			// готовим трейлеры и заголовки для получения в ответе
			var trailer metadata.MD
			callOpts := append(opts, grpc.Trailer(&trailer))

			lastErr = invoker(ctx, method, req, reply, cc, callOpts...)
			// запрос завершился без ошибки
			if lastErr == nil {
				if throttle != nil {
					// засчитываем успех и пополняем токены
					throttle.SuccessfulRPC()
				}
				successTotal.WithLabelValues(method).Inc()
				return nil
			}

			// является ли код ошибки пригодным для retry?
			if !isRetryableCode(lastErr, state.cfg.RetryableStatusCodes, state.cfg.retryableCodesSet()) {
				return lastErr
			}

			// Считаем retry для всех attempt>0, которые реально выполнились
			if attempt > 0 {
				retriesTotal.WithLabelValues(method).Inc()
			}

			// смотрим, дошли ли мы до последней попытки?
			if attempt == state.cfg.MaxAttempts-1 {
				return lastErr
			}

			duration := r.backoff(attempt, state.cfg)
			// Ждем перед повторной попыткой
			if err := wait(ctx, duration); err != nil {
				return err
			}
		}
		return nil
	}
}

func (r *Retry) backoff(attempt int, cfg MainConfig) time.Duration {
	// используем attemptsBeforePushBack для явного разделения сценариев
	// Ибо в случае с pushback - сервер явно говорит клиенту о том, что ему надо пересмотреть свою backoff стратегию
	fact := math.Pow(cfg.BackoffMultiplier, float64(attempt))
	cur := min(float64(cfg.InitialBackoff)*fact, float64(cfg.MaxBackoff))

	// Применяем рандомный jitter с фактором между 0.8 and 1.2
	cur *= 0.8 + 0.4*rand.Float64()

	return time.Duration(cur)
}

func wait(ctx context.Context, delay time.Duration) error {
	// Ждем перед повторной попыткой
	select {
	case <-ctx.Done():
		return status.FromContextError(ctx.Err()).Err()
	case <-time.After(delay):
	}
	return nil
}
