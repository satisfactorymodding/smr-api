CREATE TYPE CompatibilityState AS ENUM ('Works', 'Damaged', 'Broken');

CREATE TYPE Compatibility AS (
    note varchar,
    state CompatibilityState
);

CREATE TYPE CompatibilityInfo AS (
    EA Compatibility,
    EXP Compatibility
);

ALTER TABLE mods
    ADD COLUMN compatibility CompatibilityInfo