-- reverse: create index "modpacktarget_modpack_id_target_name" to table: "modpack_targets"
DROP INDEX "modpacktarget_modpack_id_target_name";
-- reverse: create "modpack_targets" table
DROP TABLE "modpack_targets";
-- reverse: create "modpack_tags" table
DROP TABLE "modpack_tags";
-- reverse: create index "modpackrelease_modpack_id_version" to table: "modpack_releases"
DROP INDEX "modpackrelease_modpack_id_version";
-- reverse: create "modpack_releases" table
DROP TABLE "modpack_releases";
-- reverse: create "modpack_mods" table
DROP TABLE "modpack_mods";
-- reverse: create "modpacks" table
DROP TABLE "modpacks";
