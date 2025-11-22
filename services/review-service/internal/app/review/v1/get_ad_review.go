package v1

import (
	"context"

	reviewV1 "review/internal/pkg/pb/review-service/review/v1"

	protodecimal "google.golang.org/genproto/googleapis/type/decimal"
)

func (i *Implementation) GetAdReview(ctx context.Context, req *reviewV1.GetAdReviewRequest) (*reviewV1.GetAdReviewResponse, error) {
	review, err := i.services.Review.GetAdReview(ctx, req.GetAdId())
	if err != nil {
		return nil, err
	}

	return &reviewV1.GetAdReviewResponse{
		Review: &reviewV1.GetAdReviewResponse_Review{
			AvgReview: &protodecimal.Decimal{Value: review.String()},
		},
	}, nil
}
