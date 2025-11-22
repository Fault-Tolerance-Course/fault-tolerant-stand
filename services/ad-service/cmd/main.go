package main

import (
	"context"

	"ad-service/internal"
)

func main() {
	ctx := context.Background()
	internal.New(ctx).Run(ctx)
}
