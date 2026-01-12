package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusNew Status = "new"
)

type Order struct {
	id        uuid.UUID
	status    Status
	clientID  uuid.UUID
	price     decimal.Decimal
	adID      uuid.UUID
	category  string
	createdAt time.Time
}

func NewOrder(
	clientID uuid.UUID,
	price decimal.Decimal,
	adID uuid.UUID,
	category string,
) *Order {
	id, _ := uuid.NewV7()
	return &Order{
		id:        id,
		status:    StatusNew,
		clientID:  clientID,
		price:     price,
		adID:      adID,
		category:  category,
		createdAt: time.Now(),
	}
}

func With(id uuid.UUID, opts ...Option) *Order {
	ad := &Order{id: id}

	for _, opt := range opts {
		opt(ad)
	}
	return ad
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) ClientID() uuid.UUID {
	return o.clientID
}

func (o *Order) Price() decimal.Decimal {
	return o.price
}

func (o *Order) AdID() uuid.UUID {
	return o.adID
}

func (o *Order) Category() string {
	return o.category
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}
