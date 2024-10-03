-- create "virustotal_results" table
CREATE TABLE "virustotal_results" (
    "id" character varying NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    "safe" boolean NOT NULL DEFAULT false,
    "hash" character varying NOT NULL,
    "file_name" character varying NOT NULL,
    "version_id" character varying NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "virustotal_results_versions_virustotal_results" FOREIGN KEY ("version_id") REFERENCES "versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "virustotalresult_file_name" to table: "virustotal_results"
CREATE INDEX "virustotalresult_file_name" ON "virustotal_results" ("file_name");
-- create index "virustotalresult_hash_version_id" to table: "virustotal_results"
CREATE UNIQUE INDEX "virustotalresult_hash_version_id" ON "virustotal_results" ("hash", "version_id");
-- create index "virustotalresult_safe" to table: "virustotal_results"
CREATE INDEX "virustotalresult_safe" ON "virustotal_results" ("safe");