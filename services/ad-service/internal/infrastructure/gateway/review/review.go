package review

import (
	"context"

	reviewV1 "ad-service/internal/pkg/pb/external/review-service/review/v1"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

type ExternalClient reviewV1.ReviewServiceClient

func NewExternalClient(conn grpc.ClientConnInterface) ExternalClient {
	return reviewV1.NewReviewServiceClient(conn)
}

type Client struct {
	external ExternalClient
}

func NewClient(external ExternalClient) *Client {
	return &Client{
		external: external,
	}
}

func (c *Client) GetAdReview(ctx context.Context, adID string) (decimal.Decimal, error) {
	response, err := c.external.GetAdReview(ctx, &reviewV1.GetAdReviewRequest{
		AdId: adID,
	})
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.RequireFromString(response.GetReview().GetAvgReview().GetValue()), nil
}
