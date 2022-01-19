create table if not exists announcements
(
    id varchar(14) not null
    constraint announcements_pkey primary key,
    message text not null,
    importance text not null,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone
)