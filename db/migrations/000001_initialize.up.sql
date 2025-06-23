create table message
(
    id          uuid        default gen_random_uuid()
        constraint message_pk
            primary key,
    phone       varchar(32) not null,
    content     text        not null,
    is_sent     boolean     not null default false,
    sent_at     timestamp,
    message_id  varchar(255),
    created_at  timestamp   not null default current_timestamp,
    updated_at  timestamp   not null default current_timestamp,
    deleted_at  timestamp
);
