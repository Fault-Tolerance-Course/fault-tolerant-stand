package create_order

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type CreateRequest struct {
	adID     uuid.UUID
	category string
	clientID uuid.UUID
	price    decimal.Decimal

	idempotencyKey string
	ordersVersion  int64
}

func NewCreateRequest(
	adID uuid.UUID,
	category string,
	clientID uuid.UUID,
	price decimal.Decimal,

	// закладываем поля для идемпотентности
	idempotencyKey string,
	ordersVersion int64,
) CreateRequest {
	return CreateRequest{
		adID:           adID,
		category:       category,
		clientID:       clientID,
		price:          price,
		idempotencyKey: idempotencyKey,
		ordersVersion:  ordersVersion,
	}
}

type Details struct {
	OrderID uuid.UUID
	Status  string

	OrdersVersion int64
}
