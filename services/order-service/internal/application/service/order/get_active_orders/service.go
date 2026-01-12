package get_active_orders

import (
	"context"

	"order-service/internal/domain/entity"

	"github.com/gofrs/uuid"
)

type Storage interface {
	GetActiveOrders(ctx context.Context, clientID uuid.UUID) (entity.OrderList, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetActive(ctx context.Context, clientID uuid.UUID) (entity.OrderList, error) {
	return s.storage.GetActiveOrders(ctx, clientID)
}
