package ad

import (
	"ad-service/internal/application/service"
	adV1 "ad-service/internal/pkg/pb/ad-service/ad/v1"
)

type Implementation struct {
	adV1.UnimplementedAdServiceServer
	services *service.Registry
}

func NewAdService(services *service.Registry) *Implementation {
	return &Implementation{services: services}
}
