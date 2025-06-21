-- reverse: modify "versions" table
ALTER TABLE "versions" ADD COLUMN "stability" character varying NOT NULL DEFAULT 'release';
