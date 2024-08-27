-- modify "mods" table
ALTER TABLE "mods" ADD COLUMN "toggle_network_use" boolean NOT NULL DEFAULT false, ADD COLUMN "toggle_explicit_content" boolean NOT NULL DEFAULT false;
