package payment

import (
	"context"

	commoddto "payment-service/internal/application/dto"
)

type StorageInterface interface {
	CreatePayment(ctx context.Context, order *commoddto.Payment) error
}

type Service struct {
	storage StorageInterface
}

func NewService(storage StorageInterface) *Service {
	return &Service{storage: storage}
}

func (s *Service) Create(ctx context.Context, request CreateRequest) error {
	return s.storage.CreatePayment(ctx, commoddto.NewPayment(
		request.clientID,
		request.orderID,
		request.amount,
	))
}
