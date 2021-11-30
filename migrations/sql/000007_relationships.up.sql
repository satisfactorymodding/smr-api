ALTER TABLE guides
    ADD CONSTRAINT guides_user_id_users_id FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE user_groups
    ADD CONSTRAINT user_groups_user_id_users_id FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE user_mods
    ADD CONSTRAINT user_mods_user_id_users_id FOREIGN KEY (user_id) REFERENCES users (id),
    ADD CONSTRAINT user_mods_mod_id_mods_id FOREIGN KEY (mod_id) REFERENCES mods (id);

ALTER TABLE user_sessions
    ADD CONSTRAINT user_sessions_user_id_users_id FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE versions
    ADD CONSTRAINT versions_mod_id_mods_id FOREIGN KEY (mod_id) REFERENCES mods (id);
