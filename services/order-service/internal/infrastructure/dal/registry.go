package dal

import (
	"order-service/internal/infrastructure/dal/order"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Registry struct {
	OrderRepository *order.Repository
}

func NewRegistry(pool *pgxpool.Pool) *Registry {
	return &Registry{
		OrderRepository: order.NewRepository(pool),
	}
}
