create table if not exists user_groups
(
    user_id    varchar(14) not null,
    group_id   varchar(14) not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    PRIMARY KEY (user_id, group_id)
);
