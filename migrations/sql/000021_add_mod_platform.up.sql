-- Mod Links
create table if not exists mod_archs
(
    id                      varchar(14) not null constraint mod_archs_pkey primary key,
    mod_version_arch_id     varchar(14),
    platform                varchar(16),
    key                     text,
    hash                    char(64),
    size                    bigint
);
create index if not exists idx_mod_arch_id on mod_archs (mod_version_arch_id, platform);

-- SML Links
create table if not exists sml_archs
(
    id                      varchar(14) not null constraint sml_archs_pkey primary key,
    sml_version_arch_id     varchar(14),
    platform                varchar(16),
    link                    text
);
create index if not exists idx_sml_archs_id on sml_archs (sml_version_arch_id, platform);