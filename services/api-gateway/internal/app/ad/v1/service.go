package ad

import (
	adV1 "api-gateway/internal/pkg/pb/api-gateway/ad/v1"

	externalV1 "api-gateway/internal/pkg/pb/external/ad-service/ad/v1"
)

type Implementation struct {
	adV1.UnimplementedAdServiceServer
	external externalV1.AdServiceClient
}

func NewAdService(external externalV1.AdServiceClient) *Implementation {
	return &Implementation{external: external}
}
