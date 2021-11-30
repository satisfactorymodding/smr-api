ALTER TABLE version_dependencies
    ADD COLUMN IF NOT EXISTS optional bool default false;
