-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "avatar_thumbhash";
-- reverse: modify "mods" table
ALTER TABLE "mods" DROP COLUMN "logo_thumbhash";
