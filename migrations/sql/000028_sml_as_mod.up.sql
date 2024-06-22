
-- Game engine versions
CREATE TABLE IF NOT EXISTS satisfactory_versions (
    id              varchar(14) PRIMARY KEY,
    version         int NOT NULL UNIQUE,
    engine_version  varchar(16) default '4.26'
);


INSERT INTO satisfactory_versions (id, version, engine_version)
SELECT generate_random_id(14), satisfactory_version, engine_version FROM sml_versions ON CONFLICT DO NOTHING;

-- SML mod
INSERT INTO mods (id, mod_reference, name, created_at, approved)
VALUES (
    generate_random_id(14),
    'SML',
    'Satisfactory Mod Loader',
    (SELECT created_at FROM sml_versions ORDER BY created_at LIMIT 1),
    true
)
ON CONFLICT (mod_reference) DO
    UPDATE SET name = EXCLUDED.name, created_at = EXCLUDED.created_at, approved = EXCLUDED.approved;

-- SML mod versions
INSERT INTO versions (id, mod_id, version, created_at, updated_at, deleted_at, stability, changelog, mod_reference, approved)
SELECT
    id,
    (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1),
    version, date, updated_at, deleted_at, stability, changelog,
    'SML',
    true
FROM sml_versions;

-- SML mod version targets
INSERT INTO version_targets (version_id, target_name, key)
SELECT
    (SELECT id FROM versions WHERE mod_reference = 'SML' AND version = (SELECT version FROM sml_versions WHERE sml_versions.id = version_id)),
    target_name,
    link -- Store the version link here for now, it will be replaced with a storage key by the code migration
FROM sml_version_targets;


-- SML satisfactory version
INSERT INTO version_dependencies (version_id, mod_id, condition)
SELECT
    id,
    'FactoryGame',
    '>=' || satisfactory_version
FROM sml_versions;

-- SML bootstrap version
INSERT INTO version_dependencies (version_id, mod_id, condition)
SELECT
    id,
    'bootstrap',
    '>=' || bootstrap_version
FROM sml_versions
WHERE bootstrap_version != '0.0.0' AND bootstrap_version IS NOT NULL;

-- SML devs
INSERT INTO user_mods (user_id, mod_id, role)
SELECT user_id, (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1), 'editor' FROM user_groups WHERE group_id = '3';

-- Drop table
DROP TABLE sml_version_targets;
DROP TABLE sml_versions;
