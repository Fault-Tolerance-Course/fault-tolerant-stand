package ad

import (
	"ad-service/internal/domain/entity"
	"context"

	"ad-service/internal/pkg/convert"
	adV1 "ad-service/internal/pkg/pb/ad-service/ad/v1"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) ListAds(ctx context.Context, _ *adV1.ListAdsRequest) (*adV1.ListAdsResponse, error) {
	ads, err := i.services.GetActiveAds.GeActiveAds(ctx)
	if err != nil {
		return nil, err
	}

	return &adV1.ListAdsResponse{
		Ads: lo.Map(ads, func(ad *entity.Ad, _ int) *adV1.ListAdsResponse_Ad {
			return &adV1.ListAdsResponse_Ad{
				Id:        ad.ID().String(),
				Title:     ad.Title(),
				AuthorId:  ad.AuthorID().String(),
				Price:     convert.DecimalToMoney(ad.Price()),
				Status:    string(ad.Status()),
				Category:  ad.Category(),
				CreatedAt: timestamppb.New(ad.CreatedAt()),
			}
		}),
	}, nil
}
