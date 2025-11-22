package main

import (
	"context"

	"review/internal"
)

func main() {
	ctx := context.Background()
	internal.New(ctx).Run(ctx)
}
