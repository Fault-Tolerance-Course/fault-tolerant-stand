package review

import (
	"context"

	"github.com/shopspring/decimal"
)

type ReviewProvider interface {
	GetReview(ctx context.Context, adID string) (decimal.Decimal, error)
}

type Service struct {
	review ReviewProvider
}

func NewService(review ReviewProvider) *Service {
	return &Service{review: review}
}

func (s *Service) GetAdReview(ctx context.Context, adID string) (decimal.Decimal, error) {
	// получаем отзывы и оценку
	review, err := s.review.GetReview(ctx, adID)
	if err != nil {
		return decimal.Zero, err
	}

	return review, nil
}
