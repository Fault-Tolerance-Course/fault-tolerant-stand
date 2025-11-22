package dal

import (
	"payment-service/internal/infrastructure/dal/payment"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Registry struct {
	PaymentRepository *payment.Repository
}

func NewRegistry(pool *pgxpool.Pool) *Registry {
	return &Registry{
		PaymentRepository: payment.NewRepository(pool),
	}
}
