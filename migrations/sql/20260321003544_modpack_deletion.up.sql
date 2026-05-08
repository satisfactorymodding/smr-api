-- create index "modpack_deleted_at" to table: "modpacks"
CREATE INDEX "modpack_deleted_at" ON "modpacks" ("deleted_at");
-- modify "user_modpacks" table
ALTER TABLE "user_modpacks" DROP CONSTRAINT "user_modpacks_modpack_fk", DROP CONSTRAINT "user_modpacks_user_fk", ALTER COLUMN "user_id" TYPE character varying, ALTER COLUMN "modpack_id" TYPE character varying, ALTER COLUMN "role" TYPE character varying, ADD CONSTRAINT "user_modpacks_modpacks_modpack" FOREIGN KEY ("modpack_id") REFERENCES "modpacks" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "user_modpacks_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
