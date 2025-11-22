package ad

import (
	"context"

	adV1 "ad-service/internal/pkg/pb/ad-service/ad/v1"

	"github.com/gofrs/uuid"
)

func (i *Implementation) Deactivate(ctx context.Context, req *adV1.DeactivateRequest) (*adV1.DeactivateResponse, error) {
	ad, err := i.services.Deactivation.Deactivate(ctx, uuid.FromStringOrNil(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &adV1.DeactivateResponse{
		Status: string(ad.Status()),
	}, nil
}
