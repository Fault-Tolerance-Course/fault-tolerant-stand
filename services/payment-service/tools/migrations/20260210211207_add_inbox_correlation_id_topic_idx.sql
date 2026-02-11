-- +goose Up
-- +goose StatementBegin
-- +goose NO TRANSACTION
create index concurrently idx_inbox_correlation_id_topic on inbox (correlation_id,topic); -- для быстрого селекта дубликатов сообщений
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index  idx_inbox_correlation_id_topic;
-- +goose StatementEnd