package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusPending Status = "pending"
)

type Payment struct {
	id        uuid.UUID
	status    Status
	clientID  string
	orderID   uuid.UUID
	amount    decimal.Decimal
	createdAt time.Time
}

func NewPayment(
	clientID string,
	orderID uuid.UUID,
	amount decimal.Decimal,
) *Payment {
	return &Payment{
		id:        uuid.New(),
		status:    StatusPending,
		clientID:  clientID,
		orderID:   orderID,
		amount:    amount,
		createdAt: time.Now(),
	}
}

func (o *Payment) ID() uuid.UUID {
	return o.id
}

func (o *Payment) Status() Status {
	return o.status
}

func (o *Payment) ClientID() string {
	return o.clientID
}

func (o *Payment) Amount() decimal.Decimal {
	return o.amount
}

func (o *Payment) OrderID() uuid.UUID {
	return o.orderID
}

func (o *Payment) CreatedAt() time.Time {
	return o.createdAt
}
