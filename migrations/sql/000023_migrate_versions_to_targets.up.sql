-- Mod Targets --
INSERT INTO version_targets (version_id, target_name, key, hash, size)
SELECT id, 'Windows', key, hash, size
FROM versions
WHERE NOT EXISTS(SELECT 1 FROM version_targets WHERE version_targets.version_id = versions.id)
ON CONFLICT DO NOTHING;

-- SML Targets --
INSERT INTO sml_version_targets (version_id, target_name, link)
SELECT id, 'Windows', replace(link, '/tag/', '/download/') || '/SML.zip'
FROM sml_versions
WHERE version LIKE '3%'
ON CONFLICT DO NOTHING;