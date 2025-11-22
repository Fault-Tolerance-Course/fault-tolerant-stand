package review

import (
	"context"

	"github.com/shopspring/decimal"
)

var reviewMap = map[string]decimal.Decimal{
	"3f97b21c-fbe5-4d93-884d-24736730849d": decimal.NewFromFloat(4.9),
	"df42a93b-dba1-4599-add3-f618202a6683": decimal.NewFromFloat(5.0),
	"fc798f40-5344-4c8d-81ce-cc668deb6b8c": decimal.NewFromFloat(1.1),
	"ebae49eb-27ac-4f54-8ee4-949a26ec33d8": decimal.NewFromFloat(0.0),
}

type Storage struct {
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) GetReview(_ context.Context, adID string) (decimal.Decimal, error) {
	return reviewMap[adID], nil
}
