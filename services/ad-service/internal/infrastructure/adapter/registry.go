package adapter

import (
	"ad-service/internal/infrastructure/adapter/ad"
	"ad-service/internal/infrastructure/adapter/ad/payload"
	"ad-service/internal/infrastructure/storage"
	"ad-service/internal/pkg/cache"

	"github.com/redis/go-redis/v9"
)

type Registry struct {
	Ad *ad.Adapter
}

func NewRegistry(redis redis.UniversalClient, dal *storage.Registry) *Registry {
	return &Registry{
		Ad: ad.NewAdapter(
			cache.NewClient[payload.Advertisement, *payload.Advertisement](redis),
			dal.Ad,
		),
	}
}
