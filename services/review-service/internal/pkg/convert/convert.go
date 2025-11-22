package convert

import (
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/money"
)

const (
	exp = 1_000_000_000
)

func MoneyToDecimal(m *money.Money) decimal.Decimal {
	if m == nil {
		return decimal.Zero
	}
	units := decimal.NewFromInt(m.Units)
	nanos := decimal.NewFromInt32(m.Nanos)
	return units.Add(nanos.Div(decimal.NewFromInt(exp)))
}

func DecimalToMoney(d decimal.Decimal) *money.Money {
	units := d.IntPart()
	nanos := d.Mod(decimal.NewFromInt(1).Mul(decimal.NewFromInt(exp))).IntPart()
	return &money.Money{
		Units: units,
		Nanos: int32(nanos),
	}
}
