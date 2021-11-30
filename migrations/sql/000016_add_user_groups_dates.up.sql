ALTER TABLE user_groups
    ADD COLUMN IF NOT EXISTS created_at timestamp with time zone,
    ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone,
    ADD COLUMN IF NOT EXISTS deleted_at timestamp with time zone;