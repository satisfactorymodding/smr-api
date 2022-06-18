-- Mod Links
create table if not exists mod_links
(
    id                varchar(14) not null constraint mod_links_pkey primary key,
    mod_version_link_id             text,
    platform          varchar(16),
    link              text,
    hash              char(64),
    size              bigint
);
create index if not exists idx_mod_links_id on mod_links (mod_version_link_id, platform);

-- SML Links
create table if not exists sml_links
(
    id                varchar(14) not null constraint sml_links_pkey primary key,
    sml_version_link_id      varchar(14),
    platform          varchar(16),
    link              text
);
create index if not exists idx_sml_links_id on sml_links (sml_version_link_id, platform);