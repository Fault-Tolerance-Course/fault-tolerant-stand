package order

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateRequest struct {
	adID     uuid.UUID
	category string
	clientID string
	price    decimal.Decimal
}

func NewCreateRequest(
	adID uuid.UUID,
	category string,
	clientID string,
	price decimal.Decimal,
) CreateRequest {
	return CreateRequest{
		adID:     adID,
		category: category,
		clientID: clientID,
		price:    price,
	}
}

type Details struct {
	OrderID uuid.UUID
	Status  string
}
