package entity

import "order-service/internal/domain/value_object"

type OrderList struct {
	Orders []*Order

	Version value_object.Version
}
