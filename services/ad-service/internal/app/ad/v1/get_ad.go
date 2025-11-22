package ad

import (
	"context"

	"ad-service/internal/pkg/convert"
	adV1 "ad-service/internal/pkg/pb/ad-service/ad/v1"

	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) GetAd(ctx context.Context, req *adV1.GetAdRequest) (*adV1.GetAdResponse, error) {
	ad, err := i.services.GetAd.GetAd(ctx, uuid.FromStringOrNil(req.GetAdId()))
	if err != nil {
		return nil, err
	}

	return &adV1.GetAdResponse{
		Ad: &adV1.GetAdResponse_Ad{
			Id:        ad.ID().String(),
			Title:     ad.Title(),
			AuthorId:  ad.AuthorID().String(),
			Price:     convert.DecimalToMoney(ad.Price()),
			Status:    string(ad.Status()),
			Category:  ad.Category(),
			CreatedAt: timestamppb.New(ad.CreatedAt()),
		},
	}, nil
}
