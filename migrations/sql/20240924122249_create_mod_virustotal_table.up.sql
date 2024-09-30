CREATE TABLE IF NOT EXISTS virustotal_results (
    id varchar(14) NOT NULL,
    hash varchar(64) NOT NULL,
    url varchar(101) NOT NULL,
    safe boolean DEFAULT false,
    version_target_id character varying(14) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    CONSTRAINT virustotal_results_pkey PRIMARY KEY (id),
    CONSTRAINT virustotal_results_hash_version_target_id_key UNIQUE (hash, version_target_id),
    CONSTRAINT virustotal_results_version_target_id_fkey FOREIGN KEY (version_target_id) REFERENCES version_targets (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION NOT VALID
);
CREATE INDEX IF NOT EXISTS virustotal_results_safe_idx ON virustotal_results USING btree (safe ASC NULLS LAST) WITH (deduplicate_items = True);
CREATE INDEX IF NOT EXISTS virustotal_results_hash_idx ON virustotal_results USING btree (hash ASC NULLS LAST) WITH (deduplicate_items = True);