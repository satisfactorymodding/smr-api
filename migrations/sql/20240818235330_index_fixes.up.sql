-- modify "announcements" table
ALTER TABLE "announcements"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
-- modify "guides" table
ALTER TABLE "guides"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
-- modify "mods" table
ALTER TABLE "mods"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
-- rename an index from "idx_mods_last_version_date" to "mod_last_version_date"
ALTER INDEX "idx_mods_last_version_date" RENAME TO "mod_last_version_date";
-- rename an index from "idx_mods_mod_reference" to "mod_mod_reference"
ALTER INDEX "idx_mods_mod_reference" RENAME TO "mod_mod_reference";
-- modify "tags" table
ALTER TABLE "tags"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
-- drop index "uix_user_groups_id" from table: "user_groups"
DROP INDEX "uix_user_groups_id";
-- modify "user_groups" table
ALTER TABLE "user_groups"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
-- drop index "user_sessions_token_key" from table: "user_sessions"
DROP INDEX "user_sessions_token_key";
-- modify "user_sessions" table
ALTER TABLE "user_sessions"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
-- rename an index from "uix_user_sessions_token" to "usersession_token"
ALTER INDEX "uix_user_sessions_token" RENAME TO "usersession_token";
-- drop index "users_email_key" from table: "users"
DROP INDEX "users_email_key";
-- modify "users" table
ALTER TABLE "users"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
-- rename an index from "idx_users_facebook_id" to "user_facebook_id"
ALTER INDEX "idx_users_facebook_id" RENAME TO "user_facebook_id";
-- rename an index from "idx_users_github_id" to "user_github_id"
ALTER INDEX "idx_users_github_id" RENAME TO "user_github_id";
-- rename an index from "idx_users_google_id" to "user_google_id"
ALTER INDEX "idx_users_google_id" RENAME TO "user_google_id";
-- rename an index from "uix_users_email" to "user_email"
ALTER INDEX "uix_users_email" RENAME TO "user_email";
-- modify "version_dependencies" table
ALTER TABLE "version_dependencies"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL;
-- drop index "uix_version_targets_id" from table: "version_targets"
DROP INDEX "uix_version_targets_id";
-- modify "versions" table
ALTER TABLE "versions"
    ALTER COLUMN "created_at" SET NOT NULL,
    ALTER COLUMN "updated_at" SET NOT NULL,
    ALTER COLUMN "game_version" SET NOT NULL;
-- rename an index from "idx_versions_approved" to "version_approved"
ALTER INDEX "idx_versions_approved" RENAME TO "version_approved";
-- rename an index from "idx_versions_denied" to "version_denied"
ALTER INDEX "idx_versions_denied" RENAME TO "version_denied";
-- rename an index from "idx_versions_mod_id" to "version_mod_id"
ALTER INDEX "idx_versions_mod_id" RENAME TO "version_mod_id";
