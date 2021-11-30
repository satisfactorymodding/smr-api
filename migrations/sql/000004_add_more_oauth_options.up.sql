ALTER TABLE users
    ADD COLUMN IF NOT EXISTS github_id   varchar(16)  NULL UNIQUE,
    ADD COLUMN IF NOT EXISTS google_id   varchar(32)  NULL UNIQUE,
    ADD COLUMN IF NOT EXISTS facebook_id varchar(128) NULL UNIQUE;

create index if not exists idx_users_github_id on users (github_id);
create index if not exists idx_users_google_id on users (google_id);
create index if not exists idx_users_facebook_id on users (facebook_id);
