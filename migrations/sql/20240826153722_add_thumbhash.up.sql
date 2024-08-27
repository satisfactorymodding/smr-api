-- modify "mods" table
ALTER TABLE "mods" ADD COLUMN "logo_thumbhash" character varying NULL;
-- modify "users" table
ALTER TABLE "users" ADD COLUMN "avatar_thumbhash" character varying NULL;
