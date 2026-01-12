package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Option func(ad *Order)

func WithPrice(price decimal.Decimal) Option {
	return func(o *Order) {
		o.price = price
	}
}

func WithClientID(clientID uuid.UUID) Option {
	return func(o *Order) {
		o.clientID = clientID
	}
}

func WithCategory(category string) Option {
	return func(o *Order) {
		o.category = category
	}
}

func WithStatus(status Status) Option {
	return func(o *Order) {
		o.status = status
	}
}

func WithAdID(adID uuid.UUID) Option {
	return func(o *Order) {
		o.adID = adID
	}
}

func WithCreatedAt(createdAt time.Time) Option {
	return func(o *Order) {
		o.createdAt = createdAt
	}
}
