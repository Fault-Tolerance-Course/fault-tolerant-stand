package create_ad

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Request struct {
	title    string
	category string
	authorID uuid.UUID
	price    decimal.Decimal
}

func NewRequest(title string, category string, authorID uuid.UUID, price decimal.Decimal) Request {
	return Request{
		title:    title,
		category: category,
		authorID: authorID,
		price:    price,
	}
}
