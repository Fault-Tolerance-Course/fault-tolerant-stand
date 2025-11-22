package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Option func(ad *Ad)

func WithTitle(title string) Option {
	return func(ad *Ad) {
		ad.title = title
	}
}

func WithPrice(price decimal.Decimal) Option {
	return func(ad *Ad) {
		ad.price = price
	}
}

func WithCategory(category string) Option {
	return func(ad *Ad) {
		ad.category = category
	}
}

func WithStatus(status Status) Option {
	return func(ad *Ad) {
		ad.status = status
	}
}

func WithAuthor(author uuid.UUID) Option {
	return func(ad *Ad) {
		ad.authorID = author
	}
}

func WithCreatedAt(createdAt time.Time) Option {
	return func(ad *Ad) {
		ad.createdAt = createdAt
	}
}

func WithVersion(current int) Option {
	return func(ad *Ad) {
		ad.version = version{
			Old: current,
		}
	}
}
