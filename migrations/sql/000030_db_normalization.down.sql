ALTER TABLE "versions"
    ALTER COLUMN "downloads" SET DEFAULT NULL,
    ALTER COLUMN "hotness" SET DEFAULT NULL;

UPDATE version_dependencies
SET mod_id = (SELECT mod_reference
              FROM mods
              WHERE id = version_dependencies.mod_id
              LIMIT 1)
WHERE mod_id NOT IN (SELECT mod_reference FROM mods);
