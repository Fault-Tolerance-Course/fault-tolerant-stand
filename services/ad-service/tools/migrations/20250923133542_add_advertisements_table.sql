-- +goose Up
-- +goose StatementBegin
create table if not exists advertisement (
        id uuid primary key,
        title text not null,
        status text not null,
        category text not null,
        author_id uuid not null,
        price numeric not null,
        version integer not null default 0,
        created_at timestamptz not null default CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists advertisement;
-- +goose StatementEnd
