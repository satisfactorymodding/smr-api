UPDATE versions
SET downloads = 0
WHERE downloads IS NULL;

UPDATE versions
SET hotness = 0
WHERE hotness IS NULL;

ALTER TABLE "versions"
    ALTER COLUMN "downloads" SET DEFAULT 0,
    ALTER COLUMN "hotness" SET DEFAULT 0,
    ALTER COLUMN "stability" TYPE character varying,
    ALTER COLUMN "stability" SET NOT NULL;

UPDATE versions
SET mod_reference = (SELECT mod_reference FROM mods WHERE mod_id = "versions".mod_id LIMIT 1)
WHERE mod_reference IS NULL;

DELETE
FROM version_dependencies
WHERE (SELECT id
       FROM mods
       WHERE mod_reference = version_dependencies.mod_id
       LIMIT 1) IS NULL
  AND mod_id != 'SML';

UPDATE version_dependencies
SET mod_id = (SELECT id
              FROM mods
              WHERE mod_reference = version_dependencies.mod_id
              LIMIT 1)
WHERE mod_id NOT IN (SELECT id FROM mods);

UPDATE versions
SET game_version = '0'
WHERE game_version IS NULL;
