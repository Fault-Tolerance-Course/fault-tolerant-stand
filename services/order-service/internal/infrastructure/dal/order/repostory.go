package order

import (
	"context"

	commondto "order-service/internal/application/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateOrder(ctx context.Context, order *commondto.Order) error {
	var sql = `insert into public.order (order_id,status,client_id,price,ad_id,category,created_at)
				values ($1,$2,$3,$4,$5,$6,$7);`

	_, err := r.pool.Exec(ctx, sql,
		order.ID(),
		order.Status(),
		order.ClientID(),
		order.Price(),
		order.AdID(),
		order.Category(),
		order.CreatedAt(),
	)
	return err
}
