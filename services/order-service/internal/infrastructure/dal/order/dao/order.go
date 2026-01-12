package dao

import (
	"order-service/internal/domain/entity"
	"order-service/internal/domain/value_object"
	"order-service/internal/pkg/xo"

	"github.com/shopspring/decimal"
)

type Order struct {
	xo.Order
	OrdersVersion int64 `db:"orders_version" json:"orders_version"`
}

func (o Order) Convert() *entity.Order {
	return entity.With(o.OrderID,
		entity.WithStatus(entity.Status(o.Status)),
		entity.WithCategory(o.Category),
		entity.WithPrice(decimal.NewFromFloat(o.Price)),
		entity.WithClientID(o.ClientID),
		entity.WithAdID(o.AdID),
		entity.WithCreatedAt(o.CreatedAt),
	)
}

type Orders []Order

func (o Orders) Convert() entity.OrderList {
	var (
		orders  = make([]*entity.Order, 0, len(o))
		version value_object.Version
	)

	for _, order := range o {
		version.Old = order.OrdersVersion

		orders = append(orders, order.Convert())
	}

	return entity.OrderList{
		Orders:  orders,
		Version: version,
	}
}
