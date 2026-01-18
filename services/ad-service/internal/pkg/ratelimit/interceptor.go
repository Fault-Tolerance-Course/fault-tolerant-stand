package ratelimit

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryServerInterceptor унарный серверный перехватчик
func (l *Limiter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var resp interface{}

		err := l.doIfAllowed(ctx, info.FullMethod, func() (err error) {
			resp, err = handler(ctx, req)
			return
		})
		
		return resp, err
	}
}
