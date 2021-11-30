create table if not exists sml_versions
(
    id                   varchar(14) not null
        constraint sml_releases_pkey primary key,
    created_at           timestamp with time zone,
    updated_at           timestamp with time zone,
    deleted_at           timestamp with time zone,

    version              varchar(32) unique,
    satisfactory_version int,
    stability            version_stability,
    date                 timestamp with time zone,
    link                 text,
    changelog            text
);
