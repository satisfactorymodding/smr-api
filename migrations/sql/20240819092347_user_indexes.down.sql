-- reverse: create index "user_google_id" to table: "users"
DROP INDEX "user_google_id";
-- reverse: create index "user_github_id" to table: "users"
DROP INDEX "user_github_id";
-- reverse: create index "user_facebook_id" to table: "users"
DROP INDEX "user_facebook_id";
-- reverse: drop index "user_google_id" from table: "users"
CREATE INDEX "user_google_id" ON "users" ("google_id");
-- reverse: drop index "user_github_id" from table: "users"
CREATE INDEX "user_github_id" ON "users" ("github_id");
-- reverse: drop index "user_facebook_id" from table: "users"
CREATE INDEX "user_facebook_id" ON "users" ("facebook_id");
