-- reverse: modify "user_modpacks" table
ALTER TABLE "user_modpacks" DROP CONSTRAINT "user_modpacks_users_user", DROP CONSTRAINT "user_modpacks_modpacks_modpack", ALTER COLUMN "role" TYPE text, ALTER COLUMN "modpack_id" TYPE character varying(14), ALTER COLUMN "user_id" TYPE character varying(14), ADD CONSTRAINT "user_modpacks_user_fk" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "user_modpacks_modpack_fk" FOREIGN KEY ("modpack_id") REFERENCES "modpacks" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: create index "modpack_deleted_at" to table: "modpacks"
DROP INDEX "modpack_deleted_at";
