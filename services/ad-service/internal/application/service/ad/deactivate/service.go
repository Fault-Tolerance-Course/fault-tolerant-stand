package deactivate

import (
	"context"

	"ad-service/internal/domain/entity"

	"github.com/gofrs/uuid"
)

type StatusChanger interface {
	UpdateStatus(ctx context.Context, ad *entity.Ad) error
}

type Getter interface {
	GetAd(ctx context.Context, adID uuid.UUID) (*entity.Ad, error)
}

type Service struct {
	statusChanged StatusChanger
	adGetter      Getter
}

func NewService(statusChanged StatusChanger, adGetter Getter) *Service {
	return &Service{statusChanged: statusChanged, adGetter: adGetter}
}

func (s *Service) Deactivate(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	// получаем объявление
	ad, err := s.adGetter.GetAd(ctx, adID)
	if err != nil {
		return nil, err
	}

	// деактивируем объявление
	ad.Deactivate()

	return ad, s.statusChanged.UpdateStatus(ctx, ad)
}
