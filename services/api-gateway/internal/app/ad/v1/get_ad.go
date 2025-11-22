package ad

import (
	"context"

	adV1 "api-gateway/internal/pkg/pb/api-gateway/ad/v1"
	externalV1 "api-gateway/internal/pkg/pb/external/ad-service/ad/v1"
)

func (i *Implementation) GetAd(ctx context.Context, req *adV1.GetAdRequest) (*adV1.GetAdResponse, error) {
	response, err := i.external.GetAd(ctx, &externalV1.GetAdRequest{
		AdId: req.GetAdId(),
	})
	if err != nil {
		return nil, err
	}

	return &adV1.GetAdResponse{
		Ad: &adV1.GetAdResponse_Ad{
			Id:        response.GetAd().GetId(),
			Title:     response.GetAd().GetTitle(),
			AuthorId:  response.GetAd().GetAuthorId(),
			Price:     response.GetAd().GetPrice(),
			Status:    response.GetAd().GetStatus(),
			Category:  response.GetAd().GetCategory(),
			CreatedAt: response.GetAd().GetCreatedAt(),
			Review: &adV1.GetAdResponse_Ad_AdReview{
				AvgReview: response.GetAd().GetReview().GetAvgReview(),
			},
		},
	}, nil
}
