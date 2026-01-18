package intercept

import (
	"review/internal/pkg/grpc/metadata/clientname"
	"review/internal/pkg/rerror"

	"context"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		internal, ok := lo.ErrorsAs[*rerror.Error](err)
		if !ok {
			return resp, err
		}

		message := internal.Error()
		switch {
		case rerror.IsCode(internal, rerror.NotFound):
			return resp, status.Error(codes.NotFound, message)
		case rerror.IsCode(internal, rerror.InvalidArgument):
			return resp, status.Error(codes.InvalidArgument, message)
		case rerror.IsCode(internal, rerror.FailedPrecondition):
			return resp, status.Error(codes.FailedPrecondition, message)
		case rerror.IsCode(internal, rerror.Conflict):
			return resp, status.Error(codes.Aborted, message)
		case rerror.IsCode(internal, rerror.External):
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
