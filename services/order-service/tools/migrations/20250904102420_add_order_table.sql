-- +goose Up
-- +goose StatementBegin
create table if not exists "order"(
    order_id uuid primary key,
    status text not null,
    client_id uuid not null,
    price numeric not null default 0,
    ad_id uuid not null,
    category text not null,
    created_at timestamptz not null default CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "order";
-- +goose StatementEnd
