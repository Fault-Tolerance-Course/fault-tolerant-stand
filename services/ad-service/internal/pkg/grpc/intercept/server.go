package intercept

import (
	"context"

	"ad-service/internal/pkg/aderror"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		internal, ok := lo.ErrorsAs[*aderror.Error](err)
		if !ok {
			return resp, err
		}

		message := internal.Error()
		switch {
		case aderror.IsCode(internal, aderror.NotFound):
			return resp, status.Error(codes.NotFound, message)
		case aderror.IsCode(internal, aderror.InvalidArgument):
			return resp, status.Error(codes.InvalidArgument, message)
		case aderror.IsCode(internal, aderror.FailedPrecondition):
			return resp, status.Error(codes.FailedPrecondition, message)
		case aderror.IsCode(internal, aderror.Conflict):
			return resp, status.Error(codes.Aborted, message)
		case aderror.IsCode(internal, aderror.External):
			return resp, status.Error(codes.Internal, message)
		}

		return resp, err
	}
}
