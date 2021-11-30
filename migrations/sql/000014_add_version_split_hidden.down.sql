ALTER TABLE versions
    DROP COLUMN version_major,
    DROP COLUMN version_minor,
    DROP COLUMN version_patch,
    DROP COLUMN size,
    DROP COLUMN hash;

ALTER TABLE mods
    DROP COLUMN hidden;

alter table version_dependencies
    alter column mod_id type varchar(14) using mod_id::varchar(14);
