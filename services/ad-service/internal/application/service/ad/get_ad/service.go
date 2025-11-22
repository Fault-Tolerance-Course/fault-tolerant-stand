package get_ad

import (
	"context"

	"ad-service/internal/domain/entity"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type AdProvider interface {
	GetAd(ctx context.Context, adID uuid.UUID) (*entity.Ad, error)
}

type ReviewProvider interface {
	GetAdReview(ctx context.Context, adID string) (decimal.Decimal, error)
}

type Service struct {
	ad     AdProvider
	review ReviewProvider
}

func NewService(ad AdProvider, review ReviewProvider) *Service {
	return &Service{ad: ad, review: review}
}

func (s *Service) GetAd(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	// получаем объявление
	ad, err := s.ad.GetAd(ctx, adID)
	if err != nil {
		return nil, err
	}

	// получаем отзывы и оценку
	review, err := s.review.GetAdReview(ctx, adID.String())
	if err != nil {
		return nil, err
	}

	ad.SetReview(review)
	return ad, nil
}
