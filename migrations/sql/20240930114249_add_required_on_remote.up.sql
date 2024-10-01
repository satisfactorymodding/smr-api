-- modify "versions" table
ALTER TABLE "versions" ADD COLUMN "required_on_remote" boolean NOT NULL DEFAULT true;
