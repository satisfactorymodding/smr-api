-- Mod Targets --
INSERT INTO version_targets (version_id, target_name, key, hash, size)
SELECT id, 'Windows', key, hash, size
FROM versions;

-- SML Targets --
INSERT INTO sml_version_targets (version_id, target_name, link)
SELECT id, 'Windows', link
FROM sml_versions
WHERE version LIKE '3%';