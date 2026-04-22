#!/bin/bash

CONTAINER="smr-api-postgres-1"

echo "Wiping existing data for a clean slate..."
docker exec -i $CONTAINER psql -U postgres -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

cat ./.mise/tasks/staging-dump.sql | docker exec -i $CONTAINER psql -U postgres

echo -e "\nImport process complete! Run "mise run api"."

echo "Press any key to continue..."
read -n 1 -s