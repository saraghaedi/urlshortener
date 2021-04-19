create table if not exists urls
(
    id              bigserial,
    url             text      not null,
    created_at      timestamp not null default now(),
    updated_at      timestamp,
    deleted_at      timestamp,
    primary key (id)
);

