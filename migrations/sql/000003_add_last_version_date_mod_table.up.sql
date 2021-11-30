ALTER TABLE mods
    ADD COLUMN IF NOT EXISTS last_version_date timestamp with time zone;

create index if not exists idx_mods_last_version_date on mods (last_version_date);

UPDATE mods
SET last_version_date = (
    SELECT created_at FROM versions WHERE approved = true AND mod_id = mods.id ORDER BY created_at DESC LIMIT 1
)
WHERE approved = true;