-- SML devs
INSERT INTO user_groups (user_id, group_id)
SELECT user_id, '3' FROM user_mods WHERE mod_id = (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1);

DELETE FROM user_mods WHERE mod_id = (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1);

-- SML versions
CREATE TABLE sml_versions
(
    id                   varchar(14) not null
        constraint sml_releases_pkey
            primary key,
    created_at           timestamp with time zone,
    updated_at           timestamp with time zone,
    deleted_at           timestamp with time zone,
    version              varchar(32) unique,
    satisfactory_version integer,
    stability            version_stability,
    date                 timestamp with time zone,
    link                 text,
    changelog            text,
    bootstrap_version    varchar(14),
    engine_version       varchar(16) default '4.26'::character varying
);

INSERT INTO sml_versions (id, version, link, satisfactory_version, bootstrap_version, date, created_at, updated_at, deleted_at, stability, changelog)
SELECT
    id,
    version,
    'https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/tag/v' || version,
    (SELECT substring(condition, 3) FROM version_dependencies WHERE version_id = id AND mod_id = 'FactoryGame' LIMIT 1)::int,
    (SELECT substring(condition, 3) FROM version_dependencies WHERE version_id = id AND mod_id = 'bootstrap' LIMIT 1),
    created_at, created_at, updated_at, deleted_at, stability, changelog
FROM versions
WHERE mod_id = (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1);

DELETE FROM version_dependencies WHERE version_id IN (SELECT id FROM versions WHERE mod_id = (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1));

-- SML version targets
CREATE TABLE sml_version_targets
(
    version_id     varchar(14) REFERENCES sml_versions (id),
    target_name    varchar(16),
    link           text
);

ALTER TABLE sml_version_targets
    ADD CONSTRAINT sml_version_targets_pkey PRIMARY KEY (version_id, target_name);

INSERT INTO sml_version_targets (version_id, target_name, link)
SELECT
    version_id,
    target_name,
    'https://github.com/satisfactorymodding/SatisfactoryModLoader/releases/download/v'
        || (SELECT version from versions WHERE id = version_id LIMIT 1)
        || '/SML'
        || CASE WHEN (SELECT COUNT(*) FROM version_targets AS tmp WHERE tmp.version_id = version_targets.version_id) > 1 THEN '-' || target_name ELSE '' END -- Append target name if there are multiple targets
        || '.zip'
FROM version_targets
WHERE version_id IN (SELECT id FROM versions WHERE mod_id = (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1));

DELETE FROM version_targets WHERE version_id IN (SELECT id FROM versions WHERE mod_id = (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1));

-- SML versions
DELETE FROM versions WHERE mod_id = (SELECT id FROM mods WHERE mod_reference = 'SML' LIMIT 1);

-- SML mod
DELETE FROM mods WHERE mod_reference = 'SML';

-- Game engine versions
UPDATE sml_versions SET engine_version = (
    SELECT engine_version
    FROM satisfactory_versions
    -- latest satisfactory version <= compiled satisfactory version
    WHERE satisfactory_versions.version <= sml_versions.satisfactory_version
    ORDER BY satisfactory_versions.version DESC
    LIMIT 1
) WHERE engine_version = '4.26';

DROP TABLE satisfactory_versions;
