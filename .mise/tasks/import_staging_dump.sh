#!/bin/bash

CONTAINER="smr-api-postgres-1"
URL="https://storage.ficsit.app/file/smr-db-dumps/staging-dump.sql"

echo -e "WARNING: This will wipe your local database and import a fresh staging dump."
read -p "Are you sure you want to proceed? (y/N): " confirm

if [[ $confirm != [yY] && $confirm != [yY][eE][sS] ]]; then
    echo "Import cancelled."
    exit 0
fi

echo "---"
echo "Wiping existing data for a clean slate..."
docker exec -i $CONTAINER psql -U postgres -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

echo "Downloading and importing staging dump from Ficsit Storage..."
# Using curl to stream the SQL directly into psql
curl -sSL "$URL" | docker exec -i $CONTAINER psql -U postgres

echo -e "Import complete! Run 'mise run api'."

echo "Press any key to exit..."
read -n 1 -s