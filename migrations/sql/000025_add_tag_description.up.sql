ALTER TABLE tags
    ADD COLUMN IF NOT EXISTS description varchar(512);