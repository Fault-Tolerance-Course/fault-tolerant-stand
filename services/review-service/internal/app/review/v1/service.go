package v1

import (
	"review/internal/application/service"
	reviewV1 "review/internal/pkg/pb/review-service/review/v1"
)

type Implementation struct {
	reviewV1.UnimplementedReviewServiceServer
	services *service.Registry

	requestCounter uint64 // atomic counter

	unstableReqCounter uint64 // Счетчик запросов для UnstableMethod
}

func NewReviewService(services *service.Registry) *Implementation {
	return &Implementation{services: services}
}
