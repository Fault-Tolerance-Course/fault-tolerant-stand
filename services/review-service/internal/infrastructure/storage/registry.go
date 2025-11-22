package storage

import "review/internal/infrastructure/storage/review"

type Registry struct {
	Review *review.Storage
}

func NewRegistry() *Registry {
	return &Registry{
		Review: review.NewStorage(),
	}
}
