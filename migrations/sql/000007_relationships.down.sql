ALTER TABLE guides
    DROP CONSTRAINT guides_user_id_users_id;

ALTER TABLE user_groups
    DROP CONSTRAINT user_groups_user_id_users_id;

ALTER TABLE user_mods
    DROP CONSTRAINT user_mods_user_id_users_id,
    DROP CONSTRAINT user_mods_mod_id_mods_id;

ALTER TABLE user_sessions
    DROP CONSTRAINT user_sessions_user_id_users_id;

ALTER TABLE versions
    DROP CONSTRAINT versions_mod_id_mods_id;
