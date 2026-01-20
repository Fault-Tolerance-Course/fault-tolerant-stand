package ctxutil

import (
	"context"
	"time"
)

// Detach возвращает контекст, который не отслеживает отмену родительского
func Detach(ctx context.Context) context.Context {
	return detachedContext{ctx: ctx}
}

type detachedContext struct {
	ctx context.Context
}

func (d detachedContext) Deadline() (deadline time.Time, ok bool) { return time.Time{}, false }

func (d detachedContext) Done() <-chan struct{} { return nil }

func (d detachedContext) Err() error {
	return nil
}

func (d detachedContext) Value(key interface{}) interface{} { return d.ctx.Value(key) }
