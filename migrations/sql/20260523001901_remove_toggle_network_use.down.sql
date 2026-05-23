-- reverse: modify "mods" table
ALTER TABLE "mods" ADD COLUMN "toggle_network_use" boolean NOT NULL DEFAULT false;
