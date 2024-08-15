-- reverse: drop "bootstrap_versions" table
CREATE TABLE "bootstrap_versions"
(
    "id"                   character varying(14) NOT NULL,
    "created_at"           timestamptz           NULL,
    "updated_at"           timestamptz           NULL,
    "deleted_at"           timestamptz           NULL,
    "version"              character varying(32) NULL,
    "satisfactory_version" integer               NULL,
    "stability"            "version_stability"   NULL,
    "date"                 timestamptz           NULL,
    "link"                 text                  NULL,
    "changelog"            text                  NULL,
    CONSTRAINT "bootstrap_versions_version_key" UNIQUE ("version")
);
-- reverse: rename an index from "uix_version_targets_id" to "versiontarget_id"
ALTER INDEX "versiontarget_id" RENAME TO "uix_version_targets_id";
-- reverse: create index "versiontarget_version_id_target_name" to table: "version_targets"
DROP INDEX "versiontarget_version_id_target_name";
-- reverse: modify "version_targets" table
ALTER TABLE "version_targets"
    DROP CONSTRAINT "version_targets_versions_targets",
    DROP CONSTRAINT "version_targets_pkey",
    ALTER COLUMN "id" TYPE character varying(14),
    ALTER COLUMN "id" DROP NOT NULL,
    ALTER COLUMN "id" SET DEFAULT generate_random_id(14),
    ALTER COLUMN "hash" TYPE character(64),
    ALTER COLUMN "key" TYPE text,
    ALTER COLUMN "target_name" TYPE character varying(16),
    ALTER COLUMN "version_id" TYPE character varying(14),
    ADD CONSTRAINT "version_targets_version_id_fkey" FOREIGN KEY ("version_id") REFERENCES "versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD PRIMARY KEY ("version_id", "target_name");
-- reverse: create index "versiondependency_deleted_at" to table: "version_dependencies"
DROP INDEX "versiondependency_deleted_at";
-- reverse: modify "version_dependencies" table
ALTER TABLE "version_dependencies"
    DROP CONSTRAINT "version_dependencies_versions_version",
    DROP CONSTRAINT "version_dependencies_mods_mod",
    ALTER COLUMN "optional" DROP NOT NULL,
    ALTER COLUMN "optional" SET DEFAULT false,
    ALTER COLUMN "condition" TYPE character varying(64),
    ALTER COLUMN "condition" DROP NOT NULL,
    ALTER COLUMN "mod_id" TYPE character varying(32),
    ALTER COLUMN "version_id" TYPE character varying(14),
    ADD CONSTRAINT "version_dependencies_version_id_fkey" FOREIGN KEY ("version_id") REFERENCES "versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: rename an index from "idx_versions_deleted_at" to "version_deleted_at"
ALTER INDEX "version_deleted_at" RENAME TO "idx_versions_deleted_at";
-- reverse: create index "version_id" to table: "versions"
DROP INDEX "version_id";
-- reverse: modify "versions" table
ALTER TABLE "versions"
    DROP CONSTRAINT "versions_mods_versions",
    ALTER COLUMN "hash" TYPE character(64),
    ALTER COLUMN "version_patch" TYPE integer,
    ALTER COLUMN "version_minor" TYPE integer,
    ALTER COLUMN "version_major" TYPE integer,
    ALTER COLUMN "mod_reference" TYPE character varying(32),
    ALTER COLUMN "mod_reference" DROP NOT NULL,
    ALTER COLUMN "metadata" TYPE text,
    ALTER COLUMN "hotness" TYPE integer,
    ALTER COLUMN "hotness" DROP NOT NULL,
    ALTER COLUMN "key" TYPE text,
    ALTER COLUMN "downloads" TYPE integer,
    ALTER COLUMN "downloads" DROP NOT NULL,
    ALTER COLUMN "changelog" TYPE text,
    ALTER COLUMN "version" TYPE character varying(16),
    ALTER COLUMN "version" DROP NOT NULL,
    ALTER COLUMN "mod_id" TYPE text,
    ALTER COLUMN "mod_id" DROP NOT NULL,
    ALTER COLUMN "id" TYPE character varying(14),
    ADD CONSTRAINT "versions_mod_id_mods_id" FOREIGN KEY ("mod_id") REFERENCES "mods" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: rename an index from "idx_user_sessions_deleted_at" to "usersession_deleted_at"
ALTER INDEX "usersession_deleted_at" RENAME TO "idx_user_sessions_deleted_at";
-- reverse: create index "usersession_id" to table: "user_sessions"
DROP INDEX "usersession_id";
-- reverse: create index "user_sessions_token_key" to table: "user_sessions"
DROP INDEX "user_sessions_token_key";
-- reverse: modify "user_sessions" table
ALTER TABLE "user_sessions"
    DROP CONSTRAINT "user_sessions_users_sessions",
    ALTER COLUMN "user_agent" TYPE text,
    ALTER COLUMN "token" TYPE character varying(512),
    ALTER COLUMN "token" DROP NOT NULL,
    ALTER COLUMN "user_id" TYPE text,
    ALTER COLUMN "user_id" DROP NOT NULL,
    ALTER COLUMN "id" TYPE character varying(14),
    ADD CONSTRAINT "user_sessions_user_id_users_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: modify "user_mods" table
ALTER TABLE "user_mods"
    DROP CONSTRAINT "user_mods_users_user",
    DROP CONSTRAINT "user_mods_mods_mod",
    ALTER COLUMN "role" TYPE text,
    ALTER COLUMN "role" DROP NOT NULL,
    ALTER COLUMN "mod_id" TYPE character varying(14),
    ALTER COLUMN "user_id" TYPE character varying(14),
    ADD CONSTRAINT "user_mods_user_id_users_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD CONSTRAINT "user_mods_mod_id_mods_id" FOREIGN KEY ("mod_id") REFERENCES "mods" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: rename an index from "uix_user_groups_id" to "usergroup_id"
ALTER INDEX "usergroup_id" RENAME TO "uix_user_groups_id";
-- reverse: create index "usergroup_user_id_group_id" to table: "user_groups"
DROP INDEX "usergroup_user_id_group_id";
-- reverse: create index "usergroup_deleted_at" to table: "user_groups"
DROP INDEX "usergroup_deleted_at";
-- reverse: modify "user_groups" table
ALTER TABLE "user_groups"
    DROP CONSTRAINT "user_groups_users_groups",
    DROP CONSTRAINT "user_groups_pkey",
    ALTER COLUMN "id" TYPE character varying(14),
    ALTER COLUMN "id" DROP NOT NULL,
    ALTER COLUMN "id" SET DEFAULT generate_random_id(14),
    ALTER COLUMN "group_id" TYPE character varying(14),
    ALTER COLUMN "user_id" TYPE character varying(14),
    ADD CONSTRAINT "user_groups_user_id_users_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD PRIMARY KEY ("user_id", "group_id");
-- reverse: modify "mod_tags" table
ALTER TABLE "mod_tags"
    DROP CONSTRAINT "mod_tags_tags_tag",
    DROP CONSTRAINT "mod_tags_mods_mod",
    ALTER COLUMN "mod_id" TYPE character varying(14),
    ALTER COLUMN "tag_id" TYPE character varying(14),
    ADD CONSTRAINT "mod_tags_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD CONSTRAINT "mod_tags_mod_id_fkey" FOREIGN KEY ("mod_id") REFERENCES "mods" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: rename an index from "idx_mods_deleted_at" to "mod_deleted_at"
ALTER INDEX "mod_deleted_at" RENAME TO "idx_mods_deleted_at";
-- reverse: create index "mod_id" to table: "mods"
DROP INDEX "mod_id";
-- reverse: modify "mods" table
ALTER TABLE "mods"
    ALTER COLUMN "hidden" DROP NOT NULL,
    ALTER COLUMN "mod_reference" TYPE character varying(32),
    ALTER COLUMN "downloads" TYPE integer,
    ALTER COLUMN "downloads" DROP NOT NULL,
    ALTER COLUMN "downloads" DROP DEFAULT,
    ALTER COLUMN "popularity" TYPE integer,
    ALTER COLUMN "popularity" DROP NOT NULL,
    ALTER COLUMN "popularity" DROP DEFAULT,
    ALTER COLUMN "hotness" TYPE integer,
    ALTER COLUMN "hotness" DROP NOT NULL,
    ALTER COLUMN "hotness" DROP DEFAULT,
    ALTER COLUMN "views" TYPE integer,
    ALTER COLUMN "views" DROP NOT NULL,
    ALTER COLUMN "views" DROP DEFAULT,
    ALTER COLUMN "creator_id" TYPE text,
    ALTER COLUMN "creator_id" DROP NOT NULL,
    ALTER COLUMN "source_url" TYPE text,
    ALTER COLUMN "logo" TYPE text,
    ALTER COLUMN "logo" DROP NOT NULL,
    ALTER COLUMN "full_description" TYPE text,
    ALTER COLUMN "full_description" DROP NOT NULL,
    ALTER COLUMN "short_description" TYPE character varying(128),
    ALTER COLUMN "short_description" DROP NOT NULL,
    ALTER COLUMN "name" TYPE character varying(32),
    ALTER COLUMN "name" DROP NOT NULL,
    ALTER COLUMN "id" TYPE character varying(14);
-- reverse: modify "guide_tags" table
ALTER TABLE "guide_tags"
    DROP CONSTRAINT "guide_tags_tags_tag",
    DROP CONSTRAINT "guide_tags_guides_guide",
    ALTER COLUMN "guide_id" TYPE character varying(14),
    ALTER COLUMN "tag_id" TYPE character varying(14),
    ADD CONSTRAINT "guide_tags_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD CONSTRAINT "guide_tags_guide_id_fkey" FOREIGN KEY ("guide_id") REFERENCES "guides" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: rename an index from "idx_tags_deleted_at" to "tag_deleted_at"
ALTER INDEX "tag_deleted_at" RENAME TO "idx_tags_deleted_at";
-- reverse: create index "tag_id" to table: "tags"
DROP INDEX "tag_id";
-- reverse: create index "tags_name_key" to table: "tags"
DROP INDEX "tags_name_key";
-- reverse: modify "tags" table
ALTER TABLE "tags"
    ALTER COLUMN "description" TYPE character varying(512),
    ALTER COLUMN "name" TYPE character varying(24),
    ALTER COLUMN "id" TYPE character varying(14),
    ADD CONSTRAINT "tags_name_key" UNIQUE ("name");
-- reverse: rename an index from "idx_guides_deleted_at" to "guide_deleted_at"
ALTER INDEX "guide_deleted_at" RENAME TO "idx_guides_deleted_at";
-- reverse: create index "guide_id" to table: "guides"
DROP INDEX "guide_id";
-- reverse: modify "guides" table
ALTER TABLE "guides"
    DROP CONSTRAINT "guides_users_guides",
    ALTER COLUMN "user_id" TYPE text,
    ALTER COLUMN "views" TYPE integer,
    ALTER COLUMN "views" DROP NOT NULL,
    ALTER COLUMN "views" DROP DEFAULT,
    ALTER COLUMN "guide" TYPE text,
    ALTER COLUMN "guide" DROP NOT NULL,
    ALTER COLUMN "short_description" TYPE character varying(128),
    ALTER COLUMN "short_description" DROP NOT NULL,
    ALTER COLUMN "name" TYPE character varying(32),
    ALTER COLUMN "name" DROP NOT NULL,
    ALTER COLUMN "id" TYPE character varying(14),
    ADD CONSTRAINT "guides_user_id_users_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: rename an index from "idx_users_deleted_at" to "user_deleted_at"
ALTER INDEX "user_deleted_at" RENAME TO "idx_users_deleted_at";
-- reverse: create index "users_email_key" to table: "users"
DROP INDEX "users_email_key";
-- reverse: create index "user_id" to table: "users"
DROP INDEX "user_id";
-- reverse: modify "users" table
ALTER TABLE "users"
    ALTER COLUMN "facebook_id" TYPE character varying(128),
    ALTER COLUMN "google_id" TYPE character varying(32),
    ALTER COLUMN "github_id" TYPE character varying(16),
    ALTER COLUMN "rank" TYPE integer,
    ALTER COLUMN "joined_from" TYPE text,
    ALTER COLUMN "avatar" TYPE text,
    ALTER COLUMN "username" TYPE character varying(32),
    ALTER COLUMN "username" DROP NOT NULL,
    ALTER COLUMN "email" TYPE character varying(256),
    ALTER COLUMN "email" DROP NOT NULL,
    ALTER COLUMN "id" TYPE character varying(14),
    ADD CONSTRAINT "users_google_id_key" UNIQUE ("google_id"),
    ADD CONSTRAINT "users_github_id_key" UNIQUE ("github_id"),
    ADD CONSTRAINT "users_facebook_id_key" UNIQUE ("facebook_id");
-- reverse: create index "announcement_id" to table: "announcements"
DROP INDEX "announcement_id";
-- reverse: create index "announcement_deleted_at" to table: "announcements"
DROP INDEX "announcement_deleted_at";
-- reverse: modify "announcements" table
ALTER TABLE "announcements"
    ALTER COLUMN "importance" TYPE text,
    ALTER COLUMN "message" TYPE text,
    ALTER COLUMN "id" TYPE character varying(14);
-- reverse: create index "satisfactoryversion_id" to table: "satisfactory_versions"
DROP INDEX "satisfactoryversion_id";
-- reverse: create index "satisfactory_versions_version_key" to table: "satisfactory_versions"
DROP INDEX "satisfactory_versions_version_key";
-- reverse: modify "satisfactory_versions" table
ALTER TABLE "satisfactory_versions"
    ALTER COLUMN "engine_version" TYPE character varying(16),
    ALTER COLUMN "engine_version" DROP NOT NULL,
    ALTER COLUMN "version" TYPE integer,
    ALTER COLUMN "id" TYPE character varying(14),
    ADD CONSTRAINT "satisfactory_versions_version_key" UNIQUE ("version");
