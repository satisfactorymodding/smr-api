-- create "modpacks" table
CREATE TABLE "modpacks"
(
    "id"                character varying NOT NULL,
    "created_at"        timestamptz       NOT NULL,
    "updated_at"        timestamptz       NOT NULL,
    "name"              character varying NOT NULL,
    "short_description" character varying NOT NULL,
    "full_description"  character varying NOT NULL,
    "logo"              character varying NULL,
    "logo_thumbhash"    character varying NULL,
    "creator_id"        character varying NOT NULL,
    "views"             bigint            NOT NULL DEFAULT 0,
    "hotness"           bigint            NOT NULL DEFAULT 0,
    "installs"          bigint            NOT NULL DEFAULT 0,
    "popularity"        bigint            NOT NULL DEFAULT 0,
    "hidden"            boolean           NOT NULL DEFAULT false,
    "parent_id"         character varying NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "modpacks_modpacks_children" FOREIGN KEY ("parent_id") REFERENCES "modpacks" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);
-- create "modpack_mods" table
CREATE TABLE "modpack_mods"
(
    "version_constraint" character varying NOT NULL,
    "modpack_id"         character varying NOT NULL,
    "mod_id"             character varying NOT NULL,
    PRIMARY KEY ("modpack_id", "mod_id"),
    CONSTRAINT "modpack_mods_modpacks_modpack" FOREIGN KEY ("modpack_id") REFERENCES "modpacks" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "modpack_mods_mods_mod" FOREIGN KEY ("mod_id") REFERENCES "mods" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create "modpack_releases" table
CREATE TABLE "modpack_releases"
(
    "id"         character varying NOT NULL,
    "created_at" timestamptz       NOT NULL,
    "updated_at" timestamptz       NOT NULL,
    "version"    character varying NOT NULL,
    "changelog"  character varying NOT NULL,
    "lockfile"   character varying NOT NULL,
    "modpack_id" character varying NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "modpack_releases_modpacks_releases" FOREIGN KEY ("modpack_id") REFERENCES "modpacks" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "modpackrelease_modpack_id_version" to table: "modpack_releases"
CREATE UNIQUE INDEX "modpackrelease_modpack_id_version" ON "modpack_releases" ("modpack_id", "version");
-- create "modpack_tags" table
CREATE TABLE "modpack_tags"
(
    "modpack_id" character varying NOT NULL,
    "tag_id"     character varying NOT NULL,
    PRIMARY KEY ("modpack_id", "tag_id"),
    CONSTRAINT "modpack_tags_modpacks_modpack" FOREIGN KEY ("modpack_id") REFERENCES "modpacks" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "modpack_tags_tags_tag" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create "modpack_targets" table
CREATE TABLE "modpack_targets"
(
    "id"          character varying NOT NULL,
    "target_name" character varying NOT NULL,
    "modpack_id"  character varying NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "modpack_targets_modpacks_targets" FOREIGN KEY ("modpack_id") REFERENCES "modpacks" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "modpacktarget_modpack_id_target_name" to table: "modpack_targets"
CREATE UNIQUE INDEX "modpacktarget_modpack_id_target_name" ON "modpack_targets" ("modpack_id", "target_name");
