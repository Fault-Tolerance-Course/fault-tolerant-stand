package intercept

import (
	"context"

	"order-service/internal/pkg/ordererror"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
