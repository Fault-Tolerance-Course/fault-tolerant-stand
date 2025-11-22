package platform

import (
	"context"

	"order-service/internal/pkg/event"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

type OutboxFlusher struct {
	pool *pgxpool.Pool
}

// NewOutboxFlusher конструктор
func NewOutboxFlusher(pool *pgxpool.Pool) *OutboxFlusher {
	return &OutboxFlusher{pool: pool}
}

// Flush функция сброса событий в outbox таблицу
func (o *OutboxFlusher) Flush(ctx context.Context, events event.Events) error {
	var (
		insert = `insert into outbox (key,scheme,message,headers)
					select unnest($1::bytea[]) as key,
						   unnest($2::text[]) as scheme,
					       unnest($3::bytea[]) as message,
					       unnest($4::jsonb[]) as headers;`

		keys    = lo.Map(events, func(event event.Event, _ int) []byte { return event.Key })
		schemas = lo.Map(events, func(event event.Event, _ int) string { return event.Schema })
		bodies  = lo.Map(events, func(event event.Event, _ int) []byte { return event.Body })
		headers = lo.Map(events, func(event event.Event, _ int) string { return string(event.Headers) })
	)

	_, err := o.pool.Exec(ctx, insert, keys, schemas, bodies, headers)
	return err
}
