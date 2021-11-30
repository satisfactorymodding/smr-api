ALTER TABLE mods
    ADD COLUMN mod_reference varchar(32) NULL;

UPDATE mods
SET mod_reference = id
WHERE mods.mod_reference IS NULL;

ALTER TABLE mods
    ALTER COLUMN mod_reference SET NOT NULL;

create unique index if not exists idx_mods_mod_reference on mods (mod_reference);

ALTER TABLE version_dependencies
    DROP CONSTRAINT version_dependencies_mod_id_fkey;