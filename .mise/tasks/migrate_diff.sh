#!/usr/bin/env bash

echo -n 'Migration Name: '

read -r migration_name

atlas migrate diff "$migration_name" --dir "file://migrations/sql?format=golang-migrate" --to "ent://db/schema" --dev-url "docker://postgres/16/test?search_path=public"
