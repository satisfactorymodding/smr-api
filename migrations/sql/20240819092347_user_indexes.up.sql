-- drop index "user_facebook_id" from table: "users"
DROP INDEX "user_facebook_id";
-- drop index "user_github_id" from table: "users"
DROP INDEX "user_github_id";
-- drop index "user_google_id" from table: "users"
DROP INDEX "user_google_id";
-- create index "user_facebook_id" to table: "users"
CREATE UNIQUE INDEX "user_facebook_id" ON "users" ("facebook_id");
-- create index "user_github_id" to table: "users"
CREATE UNIQUE INDEX "user_github_id" ON "users" ("github_id");
-- create index "user_google_id" to table: "users"
CREATE UNIQUE INDEX "user_google_id" ON "users" ("google_id");
