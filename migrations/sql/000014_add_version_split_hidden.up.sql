ALTER TABLE versions
    ADD COLUMN version_major int      NULL,
    ADD COLUMN version_minor int      NULL,
    ADD COLUMN version_patch int      NULL,
    ADD COLUMN size          bigint   NULL,
    ADD COLUMN hash          char(64) NULL;

ALTER TABLE mods
    ADD COLUMN hidden bool default false;

alter table version_dependencies
    alter column mod_id type varchar(32) using mod_id::varchar(32);