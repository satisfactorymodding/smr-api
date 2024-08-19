-- reverse: rename an index from "idx_versions_mod_id" to "version_mod_id"
ALTER INDEX "version_mod_id" RENAME TO "idx_versions_mod_id";
-- reverse: rename an index from "idx_versions_denied" to "version_denied"
ALTER INDEX "version_denied" RENAME TO "idx_versions_denied";
-- reverse: rename an index from "idx_versions_approved" to "version_approved"
ALTER INDEX "version_approved" RENAME TO "idx_versions_approved";
-- reverse: modify "versions" table
ALTER TABLE "versions"
    ALTER COLUMN "game_version" DROP NOT NULL,
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
-- reverse: drop index "uix_version_targets_id" from table: "version_targets"
CREATE UNIQUE INDEX "uix_version_targets_id" ON "version_targets" ("id");
-- reverse: modify "version_dependencies" table
ALTER TABLE "version_dependencies"
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
-- reverse: rename an index from "uix_users_email" to "user_email"
ALTER INDEX "user_email" RENAME TO "uix_users_email";
-- reverse: rename an index from "idx_users_google_id" to "user_google_id"
ALTER INDEX "user_google_id" RENAME TO "idx_users_google_id";
-- reverse: rename an index from "idx_users_github_id" to "user_github_id"
ALTER INDEX "user_github_id" RENAME TO "idx_users_github_id";
-- reverse: rename an index from "idx_users_facebook_id" to "user_facebook_id"
ALTER INDEX "user_facebook_id" RENAME TO "idx_users_facebook_id";
-- reverse: modify "users" table
ALTER TABLE "users"
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
-- reverse: drop index "users_email_key" from table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- reverse: rename an index from "uix_user_sessions_token" to "usersession_token"
ALTER INDEX "usersession_token" RENAME TO "uix_user_sessions_token";
-- reverse: modify "user_sessions" table
ALTER TABLE "user_sessions"
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
-- reverse: drop index "user_sessions_token_key" from table: "user_sessions"
CREATE UNIQUE INDEX "user_sessions_token_key" ON "user_sessions" ("token");
-- reverse: modify "user_groups" table
ALTER TABLE "user_groups"
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
-- reverse: drop index "uix_user_groups_id" from table: "user_groups"
CREATE UNIQUE INDEX "uix_user_groups_id" ON "user_groups" ("id");
-- reverse: modify "tags" table
ALTER TABLE "tags"
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
-- reverse: rename an index from "idx_mods_mod_reference" to "mod_mod_reference"
ALTER INDEX "mod_mod_reference" RENAME TO "idx_mods_mod_reference";
-- reverse: rename an index from "idx_mods_last_version_date" to "mod_last_version_date"
ALTER INDEX "mod_last_version_date" RENAME TO "idx_mods_last_version_date";
-- reverse: modify "mods" table
ALTER TABLE "mods"
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
-- reverse: modify "guides" table
ALTER TABLE "guides"
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
-- reverse: modify "announcements" table
ALTER TABLE "announcements"
    ALTER COLUMN "updated_at" DROP NOT NULL,
    ALTER COLUMN "created_at" DROP NOT NULL;
