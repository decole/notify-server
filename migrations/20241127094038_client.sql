-- +goose Up
-- +goose StatementBegin
create table client
(
    name      varchar not null
        constraint client_pku
            primary key,
    is_active boolean default true
);

alter table client
    owner to notify;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table client;
-- +goose StatementEnd
