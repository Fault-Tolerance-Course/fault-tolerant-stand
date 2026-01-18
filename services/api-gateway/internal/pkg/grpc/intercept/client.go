package intercept

import (
	"context"

	"api-gateway/internal/pkg/grpc/clientname"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func SetClientNameInterceptor(clientName string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// создаём / дополняем MD
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		md.Set(clientname.Header, clientName)

		// пушим в контекст
		ctx = metadata.NewOutgoingContext(ctx, md)

		// вызываем сам rpc
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
