package ad

import (
	"context"

	"ad-service/internal/application/service/ad/create_ad"
	"ad-service/internal/pkg/convert"

	adV1 "ad-service/internal/pkg/pb/ad-service/ad/v1"

	"github.com/gofrs/uuid"
)

func (i *Implementation) CreateAd(ctx context.Context, req *adV1.CreateAdRequest) (*adV1.CreateAdResponse, error) {
	//select {
	//case <-ctx.Done():
	//	return nil, status.Error(codes.Canceled, "context cancelled")
	//case <-time.After(20 * time.Second):
	//}

	request := create_ad.NewRequest(
		req.GetTitle(),
		req.GetCategory(),
		uuid.FromStringOrNil(req.GetAuthorId()),
		convert.MoneyToDecimal(req.GetPrice()),
	)

	ad, err := i.services.Creation.Create(ctx, request)
	if err != nil {
		return nil, err
	}

	return &adV1.CreateAdResponse{
		Ad: &adV1.CreateAdResponse_Ad{
			Id:     ad.ID().String(),
			Status: string(ad.Status()),
		},
	}, nil
}
