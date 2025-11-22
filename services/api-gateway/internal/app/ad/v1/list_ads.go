package ad

import (
	"context"

	adV1 "api-gateway/internal/pkg/pb/api-gateway/ad/v1"
	externalV1 "api-gateway/internal/pkg/pb/external/ad-service/ad/v1"
)

func (i *Implementation) ListAds(ctx context.Context, req *adV1.ListAdsRequest) (*adV1.ListAdsResponse, error) {
	response, err := i.external.ListAds(ctx, &externalV1.ListAdsRequest{})
	if err != nil {
		return nil, err
	}

	var mapped = make([]*adV1.ListAdsResponse_Ad, 0, len(response.GetAds()))

	for _, ad := range response.GetAds() {
		mapped = append(mapped, &adV1.ListAdsResponse_Ad{
			Id:        ad.GetId(),
			Title:     ad.GetTitle(),
			AuthorId:  ad.GetAuthorId(),
			Price:     ad.GetPrice(),
			Status:    ad.GetStatus(),
			Category:  ad.GetCategory(),
			CreatedAt: ad.GetCreatedAt(),
		})
	}

	return &adV1.ListAdsResponse{
		Ads: mapped,
	}, nil
}
