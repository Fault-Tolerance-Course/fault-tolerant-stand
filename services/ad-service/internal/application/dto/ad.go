package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type Ad struct {
	ID        string
	Title     string
	Desc      string
	Category  string
	Author    string
	Price     decimal.Decimal
	Status    string
	CreatedAt time.Time
	Review    Review
}

func (a *Ad) SetReview(review decimal.Decimal) {
	a.Review.Avg = review
}

type Review struct {
	Avg decimal.Decimal
}
