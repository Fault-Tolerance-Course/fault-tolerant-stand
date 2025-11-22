package transaction

import (
	"context"

	"ad-service/internal/pkg/transaction/abstract"
)

func Exec(ctx context.Context, callback func(ctx context.Context) error) (err error) {
	if tx := abstract.TxFromContext(ctx); tx != nil {
		var commit abstract.CommitFunc
		ctx, commit = abstract.ContextWithTx(ctx, nil)
		defer func() {
			err = commit(err)
		}()
	}
	return callback(ctx)
}
