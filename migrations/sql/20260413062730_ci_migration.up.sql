-- modify "modpack_targets" table
ALTER TABLE "modpack_targets" DROP CONSTRAINT "modpack_targets_modpacks_targets", ADD CONSTRAINT "modpack_targets_modpack_releases_targets" FOREIGN KEY ("version_id") REFERENCES "modpack_releases" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
