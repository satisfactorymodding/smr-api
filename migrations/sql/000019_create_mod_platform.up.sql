-- Mod Links
create table if not exists mod_links
(
    id                varchar(14) not null,
    modid             text,
    versionid         text,
    platform          varchar(8),
    side              varchar(8),
    link              text,
    hash              char(64),
    size              bigint
);
create index if not exists idx_mod_links_id on mod_links (id);

-- SML Links
create table if not exists sml_links
(
    id                varchar(14) not null,
    smlid             text,
    platform          varchar(8),
    side              varchar(8),
    link              text
);
create index if not exists idx_sml_links_id on sml_links (id);