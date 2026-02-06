// retry/retry_test.go
package retry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MockInvoker struct {
	mock.Mock
}

func (m *MockInvoker) Invoke(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	opts ...grpc.CallOption,
) error {
	args := m.Called(ctx, method, req, reply, cc, opts)
	return args.Error(0)
}

func TestUnaryInterceptor_PushbackReject(t *testing.T) {
	mockInvoker := new(MockInvoker)
	ctx := context.Background()
	req := "req"
	reply := "reply"
	var cc *grpc.ClientConn

	errRetry := status.Error(codes.Unavailable, "retry")

	mockInvoker.
		On("Invoke", mock.Anything, "/svc/method", req, reply, cc, mock.Anything).
		Run(func(args mock.Arguments) {
			// args[5] → opts []grpc.CallOption
			opts := args.Get(5).([]grpc.CallOption)
			for _, opt := range opts {
				if tr, ok := opt.(grpc.TrailerCallOption); ok {
					md := metadata.Pairs("grpc-retry-pushback-ms", "-1")
					*tr.TrailerAddr = md // <-- вот сюда надо писать трейлер
				}
			}
		}).
		Return(errRetry).Once()

	cfg := Config{
		Default: DefaultConfig{
			MainConfig: MainConfig{
				Enabled:              true,
				MaxAttempts:          3,
				InitialBackoff:       time.Millisecond,
				RetryableStatusCodes: []string{"Unavailable"},
			},
		},
	}

	r := NewRetry(cfg)
	interceptor := r.UnaryClientInterceptor()

	err := interceptor(ctx, "/svc/method", req, reply, cc, mockInvoker.Invoke)

	assert.Equal(t, errRetry, err)
	mockInvoker.AssertNumberOfCalls(t, "Invoke", 1)
}

func TestUnaryInterceptor_WithThrottler(t *testing.T) {
	mockInvoker := new(MockInvoker)
	ctx := context.Background()
	req := "req"
	reply := "reply"
	var cc *grpc.ClientConn

	errRetry := status.Error(codes.Unavailable, "retry")

	// invoker всегда возвращает retryable
	mockInvoker.
		On("Invoke", ctx, "/svc/method", req, reply, cc, mock.Anything).
		Return(errRetry).Times(5)

	cfg := Config{
		Default: DefaultConfig{
			MainConfig: MainConfig{
				Enabled:              true,
				MaxAttempts:          5,
				InitialBackoff:       time.Millisecond,
				RetryableStatusCodes: []string{"Unavailable"},
				Throttler: ThrottlerConfig{
					Enabled:    true,
					MaxTokens:  2,
					TokenRatio: 0.1, // не важно
				},
			},
		},
	}

	r := NewRetry(cfg)
	interceptor := r.UnaryClientInterceptor()

	err := interceptor(ctx, "/svc/method", req, reply, cc, mockInvoker.Invoke)
	assert.Equal(t, errRetry, err)

	// ожидание: 2 успешные consume токена → 3-я попытка троттлится
	mockInvoker.AssertNumberOfCalls(t, "Invoke", 2)
}

func TestUnaryInterceptor_PushbackDelayThenSuccess(t *testing.T) {
	mockInvoker := new(MockInvoker)
	ctx := context.Background()
	req := "req"
	reply := "reply"
	var cc *grpc.ClientConn

	errRetry := status.Error(codes.Unavailable, "retry")

	// 1-й вызов → возвращает ошибку + pushback delay
	mockInvoker.
		On("Invoke", mock.Anything, "/svc/method", req, reply, cc, mock.Anything).
		Run(func(args mock.Arguments) {
			opts := args.Get(5).([]grpc.CallOption)
			for _, opt := range opts {
				if tr, ok := opt.(grpc.TrailerCallOption); ok {
					md := metadata.Pairs("grpc-retry-pushback-ms", "50")
					*tr.TrailerAddr = md // записываем pushback трейлер
				}
			}
		}).
		Return(errRetry).Once()

	// 2-й вызов → успешный
	mockInvoker.
		On("Invoke", mock.Anything, "/svc/method", req, reply, cc, mock.Anything).
		Return(nil).Once()

	cfg := Config{
		Default: DefaultConfig{
			MainConfig: MainConfig{
				Enabled:              true,
				MaxAttempts:          3,
				InitialBackoff:       time.Millisecond,
				RetryableStatusCodes: []string{"Unavailable"},
			},
		},
	}

	r := NewRetry(cfg)
	interceptor := r.UnaryClientInterceptor()

	start := time.Now()
	err := interceptor(ctx, "/svc/method", req, reply, cc, mockInvoker.Invoke)
	assert.NoError(t, err)

	elapsed := time.Since(start)
	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond, "expected at least 50ms delay due to pushback")
	mockInvoker.AssertNumberOfCalls(t, "Invoke", 2)
}
