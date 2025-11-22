-- +goose Up
-- +goose StatementBegin
create table if not exists payment (
        id  uuid primary key,
        order_id uuid not null,
        client_id uuid not null,
        status text not null,
        amount numeric not null,
        created_at timestamptz not null default CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists payment;
-- +goose StatementEnd
