package retry

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"time"

	"api-gateway/internal/pkg/retry/pushback"

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

		var (
			throttle = state.throttler

			attemptsSincePushBack int
		)

		for attempt := 0; attempt < state.cfg.MaxAttempts; attempt++ {
			slog.Info("retry: invoke", "method", method, "attempts", attempt)
			// Проверяем, отменен ли контекст
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// добавляем x-retry-attempt & x-retry-max-attempts для демонстрации
			// в production такое делать не обязательно :)
			md := metadata.Pairs(
				"x-retry-attempt", fmt.Sprintf("%d", attempt),
				"x-retry-max-attempts", fmt.Sprintf("%d", state.cfg.MaxAttempts),
			)
			ctxWithMD := metadata.NewOutgoingContext(ctx, md)

			// готовим трейлеры и заголовки для получения в ответе
			var trailer metadata.MD
			callOpts := append(opts, grpc.Trailer(&trailer))

			err := invoker(ctxWithMD, method, req, reply, cc, callOpts...)

			// запрос завершился без ошибки
			if err == nil {
				if throttle != nil {
					// засчитываем успех и пополняем токены
					throttle.SuccessfulRPC()
				}

				slog.Info("retry: succeeded", "method", method, "attempts", attempt)
				return nil
			}

			// является ли код ошибки пригодным для retry?
			if !isRetryableCode(err, state.cfg.RetryableStatusCodes, state.cfg.retryableCodesSet()) {
				slog.Warn("retry: non-retryable error", "method", method, "err", err)
				return err
			}

			// смотрим, дошли ли мы до последней попытки?
			if attempt == state.cfg.MaxAttempts-1 {
				slog.Warn("retry: max attempts reached",
					"method", method,
					"attempts", state.cfg.MaxAttempts,
					"err", err,
				)
				return err
			}

			// parse pushback
			pb := pushback.Parse(trailer)

			// если сервер передал пушбек
			if pb.Has {
				slog.Info("retry: server pushback",
					"method", method,
					"reject", pb.Reject,
					"delay", pb.Delay,
				)

				if pb.Reject {
					// сервер сказал "не ретраить"
					return err
				}

				// сервер дал delay -> сбрасываем backoff
				attemptsSincePushBack = 0

				if err = wait(ctx, pb.Delay); err != nil {
					return err
				}
				// retry после pushback -> без throttle check
				continue
			}

			// смотрим, не исчерпан ли бюджет retry
			if throttle != nil && throttle.Throttle() {
				slog.Warn("retry: throttled by client throttler", "method", method)
				return err
			}

			// если pushback не было — fallback на backoff
			duration := r.backoff(attemptsSincePushBack, state.cfg)
			attemptsSincePushBack++

			slog.Info("retry: backoff",
				"method", method,
				"delay", duration,
				"attempts_since_pb", attemptsSincePushBack,
				"attempt", attempt+1,
			)

			// Ждем перед повторной попыткой
			if err = wait(ctx, duration); err != nil {
				return err
			}
		}
		return nil
	}
}

func (r *Retry) backoff(attemptsBeforePushBack int, cfg MainConfig) time.Duration {
	// используем attemptsBeforePushBack для явного разделения сценариев
	// Ибо в случае с pushback - сервер явно говорит клиенту о том, что ему надо пересмотреть свою backoff стратегию
	fact := math.Pow(cfg.BackoffMultiplier, float64(attemptsBeforePushBack))
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
