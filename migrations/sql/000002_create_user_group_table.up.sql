-- User => Group

create table if not exists user_groups
(
    user_id  varchar(14) not null,
    group_id varchar(14) not null,
    constraint user_groups_pkey primary key (user_id, group_id)
);
