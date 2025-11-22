package main

import (
	"context"

	"payment-service/internal"
)

func main() {
	ctx := context.Background()
	internal.New(ctx).Run(ctx)
}
