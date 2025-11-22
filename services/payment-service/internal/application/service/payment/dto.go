package payment

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateRequest struct {
	orderID  uuid.UUID
	clientID string
	amount   decimal.Decimal
}

func NewCreateRequest(
	orderID uuid.UUID,
	clientID string,
	amount decimal.Decimal,
) CreateRequest {
	return CreateRequest{
		orderID:  orderID,
		clientID: clientID,
		amount:   amount,
	}
}
