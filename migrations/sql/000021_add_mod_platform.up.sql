-- Mod Links
create table if not exists mod_links
(
    id                varchar(14) not null constraint mod_links_pkey primary key,
    mod_version_link_id             text,
    platform          varchar(8),
    side              varchar(8),
    link              text,
    hash              char(64),
    size              bigint
);
create index if not exists idx_mod_links_id on mod_links (id, platform, side);

-- SML Links
create table if not exists sml_links
(
    id                varchar(14) not null constraint sml_links_pkey primary key,
    sml_version_link_id      text,
    platform          varchar(8),
    side              varchar(8),
    link              text
);
create index if not exists idx_sml_links_id on sml_links (id, platform, side);