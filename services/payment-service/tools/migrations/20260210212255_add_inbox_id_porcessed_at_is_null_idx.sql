-- +goose Up
-- +goose StatementBegin
-- +goose NO TRANSACTION
create index concurrently idx_inbox_id_processed_at_error_counter on inbox (id) include(topic) where processed_at is null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index  idx_inbox_id_processed_at_error_counter;
-- +goose StatementEnd