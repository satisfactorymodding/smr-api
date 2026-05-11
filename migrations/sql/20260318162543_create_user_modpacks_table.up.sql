CREATE TABLE IF NOT EXISTS user_modpacks (
    user_id VARCHAR(14) NOT NULL,
    modpack_id VARCHAR(14) NOT NULL,
    role TEXT NOT NULL,
    CONSTRAINT user_modpacks_pkey PRIMARY KEY (user_id, modpack_id),
    CONSTRAINT user_modpacks_user_fk FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT user_modpacks_modpack_fk FOREIGN KEY (modpack_id) REFERENCES modpacks(id) ON UPDATE NO ACTION ON DELETE NO ACTION
);