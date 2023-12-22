# SMR API [![build](https://github.com/satisfactorymodding/smr-api/actions/workflows/build.yml/badge.svg)](https://github.com/satisfactorymodding/smr-api/actions/workflows/build.yml) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/satisfactorymodding/smr-api) [![codecov](https://codecov.io/gh/satisfactorymodding/smr-api/branch/master/graph/badge.svg?token=LFNKYWS0N2)](https://codecov.io/gh/satisfactorymodding/smr-api) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/satisfactorymodding/smr-api)

The Satisfactory Mod Repository backend API

## Running

If you are under Linux, you will need to install the following packages (or your distro's equivalent):

```bash
sudo apt update && sudo apt install -y build-essential libpng-dev
```

You will also need to generate the GQL server and REST docs via:

```bash
go generate -x -tags tools ./...
```

To start the API, execute:

```bash
go run cmd/api/serve.go
```

### Configuration

Running the API has a lot of pre-requisites.

To run the API, you will need to have a working Postgres, Redis and Storage. There is a dev composefile that you can
start via:

```bash
docker-compose -f docker-compose-dev.yml up -d
```

It is suggested you create a configuration file at `config.json` (but you can also use environment variables).

Main configuration options:

1. Postgres (started with dev composefile)
2. Redis (started with dev composefile)
3. B2 or S3 (or anything S3-compatible e.g. minio (started with dev composefile))
4. GitHub OAuth (https://github.com/settings/developers)
5. Google OAuth (https://console.developers.google.com/)
6. Facebook OAuth (https://developers.facebook.com/apps/)
7. Paseto keys (generated via `go run cmd/paseto/main.go`)
8. Frontend URL (needed for Google OAuth, otherwise can be ignored)
9. VirusTotal API key (https://www.virustotal.com/gui/sign-in)

The config format can be seen in `config/config.go` (each dot means a new level of nesting).

After startup requires the following minio commands to be executed:

```shell
mc alias set local http://localhost:9000 minio minio123
mc admin user svcacct add local minio --access-key REPLACE_ME_KEY --secret-key REPLACE_ME_SECRET
mc anonymous set public local/smr
```

## Contributing

Before contributing, please run the [linter](https://golangci-lint.run/) to ensure the code is clean and well-formed:

```bash
golangci-lint run
```

For some simple formatting issues you can use the `--fix` flag, but for more complex issues you will need to fix the code:

```bash
golangci-lint run --fix
```