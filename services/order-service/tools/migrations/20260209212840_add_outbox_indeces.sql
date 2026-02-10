-- +goose Up
-- +goose StatementBegin
-- +goose NO TRANSACTION
create index concurrently outbox_id_processed_at_null_idx on public.outbox(id) where processed_at is null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index concurrently outbox_id_processed_at_null_idx;
-- +goose StatementEnd
