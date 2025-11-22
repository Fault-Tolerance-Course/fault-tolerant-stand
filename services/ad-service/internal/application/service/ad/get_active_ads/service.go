package get_active_ads

import (
	"context"

	"ad-service/internal/domain/entity"

	"github.com/gofrs/uuid"
)

type AdProvider interface {
	GetAds(ctx context.Context, ids ...uuid.UUID) (entity.Ads, error)
}

type Service struct {
	ad AdProvider
}

func NewService(ad AdProvider) *Service {
	return &Service{ad: ad}
}

func (s *Service) GeActiveAds(ctx context.Context) (entity.Ads, error) {
	// получаем объявления
	ads, err := s.ad.GetAds(ctx)
	if err != nil {
		return nil, err
	}

	// забираем активные
	return ads.Active(), nil
}
