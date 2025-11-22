package main

import (
	"context"

	"order-service/internal"
)

func main() {
	ctx := context.Background()
	internal.New(ctx).Run(ctx)
}
