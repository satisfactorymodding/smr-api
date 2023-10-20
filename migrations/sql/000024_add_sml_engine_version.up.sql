ALTER TABLE sml_versions
    ADD COLUMN IF NOT EXISTS engine_version varchar(16) default '4.26';