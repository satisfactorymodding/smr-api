#!/usr/bin/env bash
#MISE hide=true

echo -n 'Migration Name: '

read -r migration_name

atlas migrate new "$migration_name" --dir "file://migrations/sql?format=golang-migrate"
