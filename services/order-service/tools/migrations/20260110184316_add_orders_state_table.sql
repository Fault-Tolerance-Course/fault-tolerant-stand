-- +goose Up
-- +goose StatementBegin
create table if not exists "orders_state"(
    client_id uuid primary key,
    version int not null default 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "orders_state";
-- +goose StatementEnd
