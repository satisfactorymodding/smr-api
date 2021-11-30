create index if not exists idx_versions_mod_id on versions (mod_id);
create index if not exists idx_versions_approved on versions (approved);
create index if not exists idx_versions_denied on versions (denied);
