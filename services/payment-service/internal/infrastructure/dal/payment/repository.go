package payment

import (
	"context"

	commoddto "payment-service/internal/application/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreatePayment(ctx context.Context, payment *commoddto.Payment) error {
	var sql = `insert into public.payment (id,status,client_id,order_id,amount,created_at)
				values ($1,$2,$3,$4,$5,$6);`

	_, err := r.pool.Exec(ctx, sql,
		payment.ID(),
		payment.Status(),
		payment.ClientID(),
		payment.OrderID(),
		payment.Amount(),
		payment.CreatedAt(),
	)
	return err
}
