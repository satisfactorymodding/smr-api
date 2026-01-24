#MISE description="Start docker containers and setup minio access"
param(
	[switch]$ci
)

$progress = if ($ci) { "plain" } else { "auto" }

docker compose --progress $progress --file docker-compose-dev.yml up --detach --wait && timeout /t 5
mc alias set local http://localhost:9000 minio minio123
mc admin user svcacct remove local/ REPLACE_ME_KEY --dp || Write-Output "It is okay for this REMOVE command to fail if the account is not found"
mc admin user svcacct add local minio --access-key REPLACE_ME_KEY --secret-key REPLACE_ME_SECRET
mc anonymous set public local/smr
