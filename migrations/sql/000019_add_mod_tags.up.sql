create table if not exists modtags
(
    name varchar(20) not null constraint modtags_pkey primary key,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone
);

create index if not exists idx_modtags_deleted_at on modtags (deleted_at);

create table if not exists mod_modtags
(
    modtag_name varchar(20) not null references modtags(name),
    mod_id varchar(16) not null references mods(id),
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone,

    primary key (mod_id, modtag_name)
);

create index if not exists idx_mod_modtags_deleted_at on mod_modtags (deleted_at);