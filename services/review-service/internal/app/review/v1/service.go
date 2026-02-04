package v1

import (
	"review/internal/application/service"
	reviewV1 "review/internal/pkg/pb/review-service/review/v1"
)

const (
	stageFast uint32 = iota
	stageSlow500
	stageSlow1s
	stageSlow3s
	stageFastAgain
)

type Implementation struct {
	reviewV1.UnimplementedReviewServiceServer
	services *service.Registry

	unstableReqCounter uint64 // Счетчик запросов для UnstableMethod
	stage              uint32
}

func NewReviewService(services *service.Registry) *Implementation {
	return &Implementation{services: services}
}
