#!/usr/bin/env bash
#MISE description="Use Atlas to create a new manual SQL migration file (see https://atlasgo.io/versioned/new)"

echo -n 'Migration Name: '

read -r migration_name

atlas migrate new "$migration_name" --dir "file://migrations/sql?format=golang-migrate"
