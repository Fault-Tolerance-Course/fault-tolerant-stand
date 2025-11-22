package service

import (
	"review/internal/application/service/review"
	"review/internal/infrastructure/storage"
)

type Registry struct {
	Review *review.Service
}

func NewRegistry(storage *storage.Registry) *Registry {
	return &Registry{
		Review: review.NewService(storage.Review),
	}
}
