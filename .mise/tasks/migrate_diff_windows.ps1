#MISE hide=true

Write-Host -NoNewline 'Migration Name: '
$migration_name = Read-Host
Write-Host "Atlas is running..."
atlas migrate diff "$migration_name" --dir "file://migrations/sql?format=golang-migrate" --to "ent://db/schema" --dev-url "docker://postgres/16/test?search_path=public"
if (-not $?) {
	Write-Error "If the above failed due to a checksum error, don't run the command it suggests, instead run ``mise run migrate_hash`` "
}
