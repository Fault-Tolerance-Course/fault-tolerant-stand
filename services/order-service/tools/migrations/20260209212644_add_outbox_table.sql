-- +goose Up
-- +goose StatementBegin
create table  if not exists outbox(
    id             bigserial                           primary key,
    entity_id      text                                not null,
    topic          text                                not null,
    key            bytea                               not null,
    body           bytea                               not null,
    headers        jsonb,
    metadata       jsonb,
    created_at     timestamp default current_timestamp not null,
    processed_at   timestamp,
    errors_count   smallint  default 0                 not null,
    error_desc     text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists outbox;
-- +goose StatementEnd