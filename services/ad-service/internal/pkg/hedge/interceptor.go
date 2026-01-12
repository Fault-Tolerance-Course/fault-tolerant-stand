package hedge

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

// UnaryClientInterceptor унарный клиентский перехватчик
func (b *Hedger) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		state, ok := b.getByMethod(method)
		if !ok || state == nil {
			slog.Debug("no hedger found for method or disabled", "method", method)
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		// Функция для вызова оригинального RPC
		rpcFunc := func(ctx context.Context) (interface{}, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			return reply, err
		}

		// Выполняем RPC с нашей политикой
		_, err := state.strategy.Execute(ctx, rpcFunc)
		if err != nil {
			return err
		}

		return nil
	}
}
