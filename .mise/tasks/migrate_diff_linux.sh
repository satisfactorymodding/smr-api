#!/usr/bin/env bash
#MISE hide=true

echo -n 'Migration Name: '

read -r migration_name

echo -n "Atlas is running..."

atlas migrate diff "$migration_name" --dir "file://migrations/sql?format=golang-migrate" --to "ent://db/schema" --dev-url "docker://postgres/16/test?search_path=public"
