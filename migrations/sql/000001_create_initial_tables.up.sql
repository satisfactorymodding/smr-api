-- Guides

create table if not exists guides
(
    id                varchar(14) not null
        constraint guides_pkey primary key,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone,
    name              varchar(32),
    short_description varchar(128),
    guide             text,
    views             integer,
    user_id           text
);

create index if not exists idx_guides_deleted_at on guides (deleted_at);


-- Mods

create table if not exists mods
(
    id                varchar(14)           not null
        constraint mods_pkey primary key,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone,
    name              varchar(32),
    short_description varchar(128),
    full_description  text,
    logo              text,
    source_url        text,
    creator_id        text,
    approved          boolean default false not null,
    views             integer,
    hotness           integer,
    popularity        integer,
    downloads         integer,
    denied            boolean default false not null
);

create index if not exists idx_mods_deleted_at on mods (deleted_at);


-- User => Mods

create table if not exists user_mods
(
    user_id varchar(14) not null,
    mod_id  varchar(14) not null,
    role    text,
    constraint user_mods_pkey primary key (user_id, mod_id)
);


-- User Sessions

create table if not exists user_sessions
(
    id         varchar(14) not null
        constraint user_sessions_pkey primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id    text,
    token      varchar(256),
    user_agent text
);

create index if not exists idx_user_sessions_deleted_at on user_sessions (deleted_at);
create unique index if not exists uix_user_sessions_token on user_sessions (token);


-- Users

create table if not exists users
(
    id          varchar(14)           not null
        constraint users_pkey primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    email       varchar(256),
    username    varchar(32),
    avatar      text,
    joined_from text,
    banned      boolean default false not null,
    rank        integer default 1     not null
);

create index if not exists idx_users_deleted_at on users (deleted_at);
create unique index if not exists uix_users_email on users (email);


-- Versions


CREATE OR REPLACE PROCEDURE create_version_stability_type()
LANGUAGE plpgsql
AS $$
BEGIN
    IF (SELECT COUNT(*) FROM pg_type WHERE typname = 'version_stability') = 0 THEN
        CREATE TYPE version_stability AS ENUM ('alpha', 'beta', 'release');
    END IF;
END;
$$;

CALL create_version_stability_type();

DROP PROCEDURE IF EXISTS create_version_stability_type();

create table if not exists versions
(
    id          varchar(14)           not null
        constraint versions_pkey primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    mod_id      text,
    version     varchar(16),
    sml_version varchar(16),
    changelog   text,
    downloads   integer,
    key         text,
    stability   version_stability,
    approved    boolean default false not null,
    hotness     integer,
    denied      boolean default false not null
);

create index if not exists idx_versions_deleted_at on versions (deleted_at);