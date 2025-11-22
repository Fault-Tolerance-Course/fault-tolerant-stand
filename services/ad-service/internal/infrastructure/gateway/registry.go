package gateway

import (
	"ad-service/config"
	"ad-service/internal/infrastructure/gateway/review"

	"google.golang.org/grpc"
)

type Registry struct {
	Review *review.Client
}

func NewRegistry(conn map[string]grpc.ClientConnInterface) *Registry {
	return &Registry{
		Review: review.NewClient(review.NewExternalClient(conn[config.ReviewService])),
	}
}
