-- user_groups

ALTER TABLE user_groups
    ADD COLUMN IF NOT EXISTS id varchar(14) default generate_random_id(14);

create unique index if not exists uix_user_groups_id on user_groups (id);

UPDATE user_groups SET id = generate_random_id(14) WHERE true;

-- sml_version_targets

ALTER TABLE sml_version_targets
    ADD COLUMN IF NOT EXISTS id varchar(14) default generate_random_id(14);

create unique index if not exists uix_sml_version_targets_id on sml_version_targets (id);

UPDATE sml_version_targets SET id = generate_random_id(14) WHERE true;

-- version_targets

ALTER TABLE version_targets
    ADD COLUMN IF NOT EXISTS id varchar(14) default generate_random_id(14);

create unique index if not exists uix_version_targets_id on version_targets (id);

UPDATE version_targets SET id = generate_random_id(14) WHERE true;