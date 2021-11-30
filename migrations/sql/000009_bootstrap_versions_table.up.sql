create table if not exists bootstrap_versions
(
    id                   varchar(14) not null,
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

alter table sml_versions
    add bootstrap_version varchar(14) null;