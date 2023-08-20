-- Mod version targets --
ALTER TABLE mod_archs RENAME TO version_targets;

DROP INDEX idx_mod_arch_id;

ALTER TABLE version_targets
    RENAME COLUMN mod_version_arch_id TO version_id;
ALTER TABLE version_targets
    RENAME COLUMN platform TO target_name;
ALTER TABLE version_targets
    DROP COLUMN id;

ALTER TABLE version_targets
    ADD CONSTRAINT version_targets_version_id_fkey FOREIGN KEY (version_id) REFERENCES versions (id),
    ADD CONSTRAINT version_targets_pkey PRIMARY KEY (version_id, target_name);

ALTER TABLE sml_archs RENAME TO sml_version_targets;

-- SML version targets --
DROP INDEX idx_sml_archs_id;

ALTER TABLE sml_version_targets
    RENAME COLUMN sml_version_arch_id TO version_id;
ALTER TABLE sml_version_targets
    RENAME COLUMN platform TO target_name;
ALTER TABLE sml_version_targets
    DROP COLUMN id;

ALTER TABLE sml_version_targets
    ADD CONSTRAINT sml_version_targets_version_id_fkey FOREIGN KEY (version_id) REFERENCES sml_versions (id),
    ADD CONSTRAINT sml_version_targets_pkey PRIMARY KEY (version_id, target_name);