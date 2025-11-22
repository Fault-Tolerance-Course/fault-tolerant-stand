package storage

import (
	"ad-service/internal/infrastructure/storage/ad"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Registry struct {
	Ad *ad.Storage
}

func NewRegistry(pool *pgxpool.Pool) *Registry {
	return &Registry{
		Ad: ad.NewStorage(pool),
	}
}
