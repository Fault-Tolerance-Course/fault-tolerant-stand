package ad

import (
	"context"

	adV1 "api-gateway/internal/pkg/pb/api-gateway/ad/v1"
	externalV1 "api-gateway/internal/pkg/pb/external/ad-service/ad/v1"
)

func (i *Implementation) CreateAd(ctx context.Context, req *adV1.CreateAdRequest) (*adV1.CreateAdResponse, error) {
	response, err := i.external.CreateAd(ctx, &externalV1.CreateAdRequest{
		Title:    req.GetTitle(),
		Category: req.GetCategory(),
		AuthorId: req.GetAuthorId(),
		Price:    req.GetPrice(),
	})
	if err != nil {
		return nil, err
	}

	return &adV1.CreateAdResponse{
		Ad: &adV1.CreateAdResponse_Ad{
			Id:     response.GetAd().GetId(),
			Status: response.GetAd().GetStatus(),
		},
	}, nil
}
