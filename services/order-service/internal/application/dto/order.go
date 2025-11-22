package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusNew Status = "new"
)

type Order struct {
	id        uuid.UUID
	status    Status
	clientID  string
	price     decimal.Decimal
	adID      uuid.UUID
	category  string
	createdAt time.Time
}

func NewOrder(
	clientID string,
	price decimal.Decimal,
	adID uuid.UUID,
	category string,
) *Order {
	return &Order{
		id:        uuid.New(),
		status:    StatusNew,
		clientID:  clientID,
		price:     price,
		adID:      adID,
		category:  category,
		createdAt: time.Now(),
	}
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) ClientID() string {
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
