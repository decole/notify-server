-- +goose Up
-- +goose StatementBegin
create table notify
(
    id      serial
        constraint notify_pk
            primary key,
    client  varchar(250) not null,
    message text         not null
);

alter table notify
    owner to notify;

create index idx_notify_user
    on notify (client);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop sequence notify_id_seq;

drop table notify;
-- +goose StatementEnd
