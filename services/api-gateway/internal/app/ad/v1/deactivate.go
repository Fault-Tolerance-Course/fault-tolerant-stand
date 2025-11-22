package ad

import (
	"context"

	adV1 "api-gateway/internal/pkg/pb/api-gateway/ad/v1"
	externalV1 "api-gateway/internal/pkg/pb/external/ad-service/ad/v1"
)

func (i *Implementation) Deactivate(ctx context.Context, req *adV1.DeactivateRequest) (*adV1.DeactivateResponse, error) {
	response, err := i.external.Deactivate(ctx, &externalV1.DeactivateRequest{
		Id: req.GetId(),
	})

	if err != nil {
		return nil, err
	}

	return &adV1.DeactivateResponse{
		Status: response.GetStatus(),
	}, nil
}
