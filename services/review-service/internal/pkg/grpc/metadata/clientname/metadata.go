package clientname

import "context"

const Header = "x-client-name"

type ctxKey uint8

const (
	nameKey ctxKey = iota

	unknown = "unknown"
)

func NewContext(parent context.Context, name string) context.Context {
	return context.WithValue(parent, nameKey, name)
}

func FromContext(ctx context.Context) string {
	value := ctx.Value(nameKey)
	if name, ok := value.(string); ok {
		return name
	}
	return unknown
}
