-- reverse: create index "virustotalresult_safe" to table: "virustotal_results"
DROP INDEX "virustotalresult_safe";
-- reverse: create index "virustotalresult_hash_version_id" to table: "virustotal_results"
DROP INDEX "virustotalresult_hash_version_id";
-- reverse: create index "virustotalresult_file_name" to table: "virustotal_results"
DROP INDEX "virustotalresult_file_name";
-- reverse: create "virustotal_results" table
DROP TABLE "virustotal_results";
