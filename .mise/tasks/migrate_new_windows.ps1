#MISE hide=true

Write-Host -NoNewline 'Migration Name: '
$migration_name = Read-Host
atlas migrate new "$migration_name" --dir "file://migrations/sql?format=golang-migrate"
if (-not $?) {
	Write-Error "If the above failed due to a checksum error, don't run the command it suggests, instead run ``mise run migrate_hash`` "
} else {
	Write-Host "Created new migration '$migration_name', check the migrations/sql/ directory for the new files"
}
