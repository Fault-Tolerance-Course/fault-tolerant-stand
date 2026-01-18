package intercept

import (
	"context"

	"order-service/internal/pkg/grpc/metadata/clientname"
	"order-service/internal/pkg/ordererror"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const ()

func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		internal, ok := lo.ErrorsAs[*ordererror.Error](err)
		if !ok {
			return resp, err
		}

		message := internal.Error()
		switch {
		case ordererror.IsCode(internal, ordererror.NotFound):
			return resp, status.Error(codes.NotFound, message)
		case ordererror.IsCode(internal, ordererror.InvalidArgument):
			return resp, status.Error(codes.InvalidArgument, message)
		case ordererror.IsCode(internal, ordererror.FailedPrecondition):
			return resp, status.Error(codes.FailedPrecondition, message)
		case ordererror.IsCode(internal, ordererror.Conflict):
			return resp, status.Error(codes.Aborted, message)
		case ordererror.IsCode(internal, ordererror.External):
			return resp, status.Error(codes.Internal, message)
		}

		return resp, err
	}
}

func ExtractClientNameInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if val := md.Get(clientname.Header); len(val) > 0 {
				// кладём в context
				ctx = clientname.NewContext(ctx, val[0])
			}
		}

		// продолжаем выполнение
		return handler(ctx, req)
	}
}
