
-- Game engine versions
CREATE TABLE IF NOT EXISTS satisfactory_versions (
    id              varchar(14) PRIMARY KEY,
    version         int NOT NULL UNIQUE,
    engine_version  varchar(16) default '4.26'
);


INSERT INTO satisfactory_versions (id, version, engine_version)
SELECT generate_random_id(14), satisfactory_version, engine_version FROM sml_versions ON CONFLICT DO NOTHING;

-- SML mod
INSERT INTO mods (id, mod_reference, name, short_description, full_description, creator_id, logo, created_at, updated_at, last_version_date, downloads, views, hotness, popularity, approved)
VALUES (
    generate_random_id(14),
    'SML',
    'Satisfactory Mod Loader',
    'Mod loading and compatibility API for Satisfactory',
    'Mod loading and compatibility API for Satisfactory',
    '', '',
    COALESCE((SELECT created_at FROM sml_versions ORDER BY created_at LIMIT 1), now()),
    COALESCE((SELECT updated_at FROM sml_versions ORDER BY updated_at LIMIT 1), now()),
    COALESCE((SELECT created_at FROM sml_versions ORDER BY created_at DESC LIMIT 1), now()),
    0, 0, 0, 0,
    true
)
ON CONFLICT (mod_reference) DO
    UPDATE SET name = EXCLUDED.name, created_at = EXCLUDED.created_at, updated_at=EXCLUDED.updated_at, approved = EXCLUDED.approved;
-- staging has a SML mod already, handle the conflict

-- version game version
ALTER TABLE versions ADD COLUMN game_version varchar; -- allow null for now, needs to be filled in by the code migration

-- SML mod versions
INSERT INTO versions (id, mod_id, version, created_at, updated_at, deleted_at, stability, changelog, game_version, version_major, version_minor, version_patch, mod_reference, approved)
SELECT
    id,
    (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1),
    version, date, updated_at, deleted_at, COALESCE(stability, 'alpha'), changelog,
    '>=' || satisfactory_version,
    SPLIT_PART(SPLIT_PART(version, '-', 1), '.', 1)::int,
    SPLIT_PART(SPLIT_PART(version, '-', 1), '.', 2)::int,
    SPLIT_PART(SPLIT_PART(version, '-', 1), '.', 3)::int,
    'SML',
    true
FROM sml_versions;
-- we lose the bootstrap version here, but the bootstrapper data is already not available on the API anymore

-- SML mod version targets
INSERT INTO version_targets (id, version_id, target_name, key)
SELECT
    generate_random_id(14),
    (SELECT id FROM versions WHERE mod_reference = 'SML' AND version = (SELECT version FROM sml_versions WHERE sml_versions.id = version_id)),
    target_name,
    link -- Store the version link here for now, it will be replaced with a storage key by the code migration
FROM sml_version_targets;

-- SML devs
INSERT INTO user_mods (user_id, mod_id, role)
SELECT user_id, (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1), 'editor' FROM user_groups WHERE group_id = '3';
-- we don't delete the group ID, because it turns into the game version editors group

-- Drop table
DROP TABLE sml_version_targets;
DROP TABLE sml_versions;

-- SML mod dependency
INSERT INTO version_dependencies (version_id, mod_id, condition, created_at, updated_at)
SELECT id, 'SML', sml_version, created_at, updated_at FROM versions WHERE sml_version IS NOT NULL
ON CONFLICT DO NOTHING;

ALTER TABLE versions DROP COLUMN sml_version;
