package v1

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	reviewV1 "review/internal/pkg/pb/review-service/review/v1"

	protodecimal "google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetAdReview(ctx context.Context, req *reviewV1.GetAdReviewRequest) (*reviewV1.GetAdReviewResponse, error) {
	count := atomic.AddUint64(&i.requestCounter, 1)
	// каждый 3-й запрос - медленный
	if count%3 == 0 {
		slog.Info(fmt.Sprintf("GetAdReview start sleep at: %s", time.Now().String()))
		select {
		case <-ctx.Done():
			return nil, status.Error(codes.Canceled, "context canceled")
		case <-time.After(700 * time.Millisecond):
		}
		slog.Info(fmt.Sprintf("GetAdReview end sleep at: %s", time.Now().String()))
	}

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
