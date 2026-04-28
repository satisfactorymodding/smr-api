#!/usr/bin/env bash
#USAGE flag "-ci" "Use plain progress output suitable for CI environments"
#MISE hide=true

ci=false
while [[ $# -gt 0 ]]; do
	case "$1" in
		-ci)
			ci=true
			;;
		*)
			;;
	esac
	shift
done

progress=$([ "$ci" = true ] && echo "plain" || echo "auto")

docker compose --progress "$progress" --file docker-compose-dev.yml up --detach --wait && echo "Waiting for 5 seconds" && sleep 5
echo "Setting up MinIO..."
mc alias set local http://localhost:9000 minio minio123
mc admin user svcacct remove local/ REPLACE_ME_KEY --dp || echo It is okay for this REMOVE command to fail if the account is not found
mc admin user svcacct add local minio --access-key REPLACE_ME_KEY --secret-key REPLACE_ME_SECRET
mc anonymous set public local/smr
echo "Setup complete"
