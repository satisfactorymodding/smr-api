create table if not exists version_dependencies
(
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,

    version_id varchar(14) not null REFERENCES versions (id),
    mod_id     varchar(14) not null REFERENCES mods (id),

    condition  varchar(64),

    PRIMARY KEY (version_id, mod_id)
);
