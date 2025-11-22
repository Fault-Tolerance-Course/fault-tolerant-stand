package create_ad

import (
	"context"

	"ad-service/internal/domain/entity"
)

type Creator interface {
	CreateAd(ctx context.Context, ad *entity.Ad) error
}

type Service struct {
	creator Creator
}

func NewService(creator Creator) *Service {
	return &Service{creator: creator}
}

func (s *Service) Create(ctx context.Context, request Request) (*entity.Ad, error) {
	ad := entity.NewAd(request.title, request.category, request.authorID, request.price)

	err := s.creator.CreateAd(ctx, ad)
	if err != nil {
		return nil, err
	}

	return ad, nil
}
