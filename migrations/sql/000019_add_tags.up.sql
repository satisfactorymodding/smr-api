create table if not exists tags
(
    id varchar(14) not null constraint tags_pkey primary key,
    name varchar(20) not null unique,

    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create index if not exists idx_tags_deleted_at on tags (deleted_at);

create table if not exists mod_tags
(
    tag_id varchar(20) not null references tags(id),
    mod_id varchar(16) not null references mods(id),
    primary key (mod_id, tag_id),

    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone


);

create index if not exists idx_mod_tags_deleted_at on mod_tags (deleted_at);

create table if not exists guide_tags
(
    tag_id varchar(20) not null references tags(id),
    guide_id varchar(16) not null references guides(id),
    primary key (guide_id, tag_id),

    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone


);

create index if not exists idx_guide_tags_deleted_at on guide_tags (deleted_at);