-- modify "satisfactory_versions" table
ALTER TABLE "satisfactory_versions"
    DROP CONSTRAINT "satisfactory_versions_version_key",
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "version" TYPE bigint,
    ALTER COLUMN "engine_version" TYPE character varying,
    ALTER COLUMN "engine_version" SET NOT NULL;
-- create index "satisfactory_versions_version_key" to table: "satisfactory_versions"
CREATE UNIQUE INDEX "satisfactory_versions_version_key" ON "satisfactory_versions" ("version");
-- create index "satisfactoryversion_id" to table: "satisfactory_versions"
CREATE UNIQUE INDEX "satisfactoryversion_id" ON "satisfactory_versions" ("id");
-- modify "announcements" table
ALTER TABLE "announcements"
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "message" TYPE character varying,
    ALTER COLUMN "importance" TYPE character varying;
-- create index "announcement_deleted_at" to table: "announcements"
CREATE INDEX "announcement_deleted_at" ON "announcements" ("deleted_at");
-- create index "announcement_id" to table: "announcements"
CREATE UNIQUE INDEX "announcement_id" ON "announcements" ("id");
-- modify "users" table
ALTER TABLE "users"
    DROP CONSTRAINT "users_facebook_id_key",
    DROP CONSTRAINT "users_github_id_key",
    DROP CONSTRAINT "users_google_id_key",
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "email" TYPE character varying,
    ALTER COLUMN "email" SET NOT NULL,
    ALTER COLUMN "username" TYPE character varying,
    ALTER COLUMN "username" SET NOT NULL,
    ALTER COLUMN "avatar" TYPE character varying,
    ALTER COLUMN "joined_from" TYPE character varying,
    ALTER COLUMN "rank" TYPE bigint,
    ALTER COLUMN "github_id" TYPE character varying,
    ALTER COLUMN "google_id" TYPE character varying,
    ALTER COLUMN "facebook_id" TYPE character varying;
-- create index "user_id" to table: "users"
CREATE UNIQUE INDEX "user_id" ON "users" ("id");
-- create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- rename an index from "idx_users_deleted_at" to "user_deleted_at"
ALTER INDEX "idx_users_deleted_at" RENAME TO "user_deleted_at";
-- modify "guides" table
ALTER TABLE "guides"
    DROP CONSTRAINT "guides_user_id_users_id",
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "name" TYPE character varying,
    ALTER COLUMN "name" SET NOT NULL,
    ALTER COLUMN "short_description" TYPE character varying,
    ALTER COLUMN "short_description" SET NOT NULL,
    ALTER COLUMN "guide" TYPE character varying,
    ALTER COLUMN "guide" SET NOT NULL,
    ALTER COLUMN "views" TYPE bigint,
    ALTER COLUMN "views" SET NOT NULL,
    ALTER COLUMN "views" SET DEFAULT 0,
    ALTER COLUMN "user_id" TYPE character varying,
    ADD CONSTRAINT "guides_users_guides" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- create index "guide_id" to table: "guides"
CREATE UNIQUE INDEX "guide_id" ON "guides" ("id");
-- rename an index from "idx_guides_deleted_at" to "guide_deleted_at"
ALTER INDEX "idx_guides_deleted_at" RENAME TO "guide_deleted_at";
-- modify "tags" table
ALTER TABLE "tags"
    DROP CONSTRAINT "tags_name_key",
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "name" TYPE character varying,
    ALTER COLUMN "description" TYPE character varying;
-- create index "tags_name_key" to table: "tags"
CREATE UNIQUE INDEX "tags_name_key" ON "tags" ("name");
-- create index "tag_id" to table: "tags"
CREATE UNIQUE INDEX "tag_id" ON "tags" ("id");
-- rename an index from "idx_tags_deleted_at" to "tag_deleted_at"
ALTER INDEX "idx_tags_deleted_at" RENAME TO "tag_deleted_at";
-- modify "guide_tags" table
ALTER TABLE "guide_tags"
    DROP CONSTRAINT "guide_tags_guide_id_fkey",
    DROP CONSTRAINT "guide_tags_tag_id_fkey",
    ALTER COLUMN "tag_id" TYPE character varying,
    ALTER COLUMN "guide_id" TYPE character varying,
    ADD CONSTRAINT "guide_tags_guides_guide" FOREIGN KEY ("guide_id") REFERENCES "guides" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD CONSTRAINT "guide_tags_tags_tag" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "mods" table
ALTER TABLE "mods"
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "name" TYPE character varying,
    ALTER COLUMN "name" SET NOT NULL,
    ALTER COLUMN "short_description" TYPE character varying,
    ALTER COLUMN "short_description" SET NOT NULL,
    ALTER COLUMN "full_description" TYPE character varying,
    ALTER COLUMN "full_description" SET NOT NULL,
    ALTER COLUMN "logo" TYPE character varying,
    ALTER COLUMN "logo" SET NOT NULL,
    ALTER COLUMN "source_url" TYPE character varying,
    ALTER COLUMN "creator_id" TYPE character varying,
    ALTER COLUMN "creator_id" SET NOT NULL,
    ALTER COLUMN "views" TYPE bigint,
    ALTER COLUMN "views" SET NOT NULL,
    ALTER COLUMN "views" SET DEFAULT 0,
    ALTER COLUMN "hotness" TYPE bigint,
    ALTER COLUMN "hotness" SET NOT NULL,
    ALTER COLUMN "hotness" SET DEFAULT 0,
    ALTER COLUMN "popularity" TYPE bigint,
    ALTER COLUMN "popularity" SET NOT NULL,
    ALTER COLUMN "popularity" SET DEFAULT 0,
    ALTER COLUMN "downloads" TYPE bigint,
    ALTER COLUMN "downloads" SET NOT NULL,
    ALTER COLUMN "downloads" SET DEFAULT 0,
    ALTER COLUMN "mod_reference" TYPE character varying,
    ALTER COLUMN "hidden" SET NOT NULL;
-- create index "mod_id" to table: "mods"
CREATE UNIQUE INDEX "mod_id" ON "mods" ("id");
-- rename an index from "idx_mods_deleted_at" to "mod_deleted_at"
ALTER INDEX "idx_mods_deleted_at" RENAME TO "mod_deleted_at";
-- modify "mod_tags" table
ALTER TABLE "mod_tags"
    DROP CONSTRAINT "mod_tags_mod_id_fkey",
    DROP CONSTRAINT "mod_tags_tag_id_fkey",
    ALTER COLUMN "tag_id" TYPE character varying,
    ALTER COLUMN "mod_id" TYPE character varying,
    ADD CONSTRAINT "mod_tags_mods_mod" FOREIGN KEY ("mod_id") REFERENCES "mods" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD CONSTRAINT "mod_tags_tags_tag" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "user_groups" table
ALTER TABLE "user_groups"
    DROP CONSTRAINT "user_groups_pkey",
    DROP CONSTRAINT "user_groups_user_id_users_id",
    ALTER COLUMN "user_id" TYPE character varying,
    ALTER COLUMN "group_id" TYPE character varying,
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "id" SET NOT NULL,
    ALTER COLUMN "id" DROP DEFAULT,
    ADD PRIMARY KEY ("id"),
    ADD CONSTRAINT "user_groups_users_groups" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- create index "usergroup_deleted_at" to table: "user_groups"
CREATE INDEX "usergroup_deleted_at" ON "user_groups" ("deleted_at");
-- create index "usergroup_user_id_group_id" to table: "user_groups"
CREATE UNIQUE INDEX "usergroup_user_id_group_id" ON "user_groups" ("user_id", "group_id");
-- rename an index from "uix_user_groups_id" to "usergroup_id"
ALTER INDEX "uix_user_groups_id" RENAME TO "usergroup_id";
-- modify "user_mods" table
ALTER TABLE "user_mods"
    DROP CONSTRAINT "user_mods_mod_id_mods_id",
    DROP CONSTRAINT "user_mods_user_id_users_id",
    ALTER COLUMN "user_id" TYPE character varying,
    ALTER COLUMN "mod_id" TYPE character varying,
    ALTER COLUMN "role" TYPE character varying,
    ALTER COLUMN "role" SET NOT NULL,
    ADD CONSTRAINT "user_mods_mods_mod" FOREIGN KEY ("mod_id") REFERENCES "mods" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD CONSTRAINT "user_mods_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "user_sessions" table
ALTER TABLE "user_sessions"
    DROP CONSTRAINT "user_sessions_user_id_users_id",
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "user_id" TYPE character varying,
    ALTER COLUMN "user_id" SET NOT NULL,
    ALTER COLUMN "token" TYPE character varying,
    ALTER COLUMN "token" SET NOT NULL,
    ALTER COLUMN "user_agent" TYPE character varying,
    ADD CONSTRAINT "user_sessions_users_sessions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- create index "user_sessions_token_key" to table: "user_sessions"
CREATE UNIQUE INDEX "user_sessions_token_key" ON "user_sessions" ("token");
-- create index "usersession_id" to table: "user_sessions"
CREATE UNIQUE INDEX "usersession_id" ON "user_sessions" ("id");
-- rename an index from "idx_user_sessions_deleted_at" to "usersession_deleted_at"
ALTER INDEX "idx_user_sessions_deleted_at" RENAME TO "usersession_deleted_at";
-- modify "versions" table
ALTER TABLE "versions"
    DROP CONSTRAINT "versions_mod_id_mods_id",
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "mod_id" TYPE character varying,
    ALTER COLUMN "mod_id" SET NOT NULL,
    ALTER COLUMN "version" TYPE character varying,
    ALTER COLUMN "version" SET NOT NULL,
    ALTER COLUMN "changelog" TYPE character varying,
    ALTER COLUMN "downloads" TYPE bigint,
    ALTER COLUMN "downloads" SET DEFAULT 0,
    ALTER COLUMN "key" TYPE character varying,
    ALTER COLUMN "stability" TYPE character varying,
    ALTER COLUMN "stability" SET NOT NULL,
    ALTER COLUMN "hotness" TYPE bigint,
    ALTER COLUMN "hotness" SET NOT NULL,
    ALTER COLUMN "hotness" SET DEFAULT 0,
    ALTER COLUMN "metadata" TYPE character varying,
    ALTER COLUMN "mod_reference" TYPE character varying,
    ALTER COLUMN "mod_reference" SET NOT NULL,
    ALTER COLUMN "version_major" TYPE bigint,
    ALTER COLUMN "version_minor" TYPE bigint,
    ALTER COLUMN "version_patch" TYPE bigint,
    ALTER COLUMN "hash" TYPE character varying,
    ADD CONSTRAINT "versions_mods_versions" FOREIGN KEY ("mod_id") REFERENCES "mods" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- create index "version_id" to table: "versions"
CREATE UNIQUE INDEX "version_id" ON "versions" ("id");
-- rename an index from "idx_versions_deleted_at" to "version_deleted_at"
ALTER INDEX "idx_versions_deleted_at" RENAME TO "version_deleted_at";
-- modify "version_dependencies" table
ALTER TABLE "version_dependencies"
    DROP CONSTRAINT "version_dependencies_version_id_fkey",
    ALTER COLUMN "version_id" TYPE character varying,
    ALTER COLUMN "mod_id" TYPE character varying,
    ALTER COLUMN "condition" TYPE character varying,
    ALTER COLUMN "condition" SET NOT NULL,
    ALTER COLUMN "optional" SET NOT NULL,
    ALTER COLUMN "optional" DROP DEFAULT,
    ADD CONSTRAINT "version_dependencies_mods_mod" FOREIGN KEY ("mod_id") REFERENCES "mods" ("mod_reference") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD CONSTRAINT "version_dependencies_versions_version" FOREIGN KEY ("version_id") REFERENCES "versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- create index "versiondependency_deleted_at" to table: "version_dependencies"
CREATE INDEX "versiondependency_deleted_at" ON "version_dependencies" ("deleted_at");
-- modify "version_targets" table
ALTER TABLE "version_targets"
    DROP CONSTRAINT "version_targets_pkey",
    DROP CONSTRAINT "version_targets_version_id_fkey",
    ALTER COLUMN "version_id" TYPE character varying,
    ALTER COLUMN "target_name" TYPE character varying,
    ALTER COLUMN "key" TYPE character varying,
    ALTER COLUMN "hash" TYPE character varying,
    ALTER COLUMN "id" TYPE character varying,
    ALTER COLUMN "id" SET NOT NULL,
    ALTER COLUMN "id" DROP DEFAULT,
    ADD PRIMARY KEY ("id"),
    ADD CONSTRAINT "version_targets_versions_targets" FOREIGN KEY ("version_id") REFERENCES "versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- create index "versiontarget_version_id_target_name" to table: "version_targets"
CREATE UNIQUE INDEX "versiontarget_version_id_target_name" ON "version_targets" ("version_id", "target_name");
-- rename an index from "uix_version_targets_id" to "versiontarget_id"
ALTER INDEX "uix_version_targets_id" RENAME TO "versiontarget_id";
-- drop "bootstrap_versions" table
DROP TABLE "bootstrap_versions";
-- drop enum type "version_stability"
DROP TYPE "version_stability";
