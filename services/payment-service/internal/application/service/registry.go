package service

import (
	"payment-service/internal/application/service/payment"
	"payment-service/internal/infrastructure/dal"
)

type Registry struct {
	Payment *payment.Service
}

func NewRegistry(dal *dal.Registry) *Registry {
	return &Registry{
		Payment: payment.NewService(dal.PaymentRepository),
	}
}
