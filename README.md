# SMR API

The Satisfactory Mod Repository backend API

## Running

Execute:

```bash
go run cmd/api/serve.go
```

### Configuration

Running the API has a lot of pre-requisites.

It is suggested you create a configuration file at `config.json` (but you can also use environment variables).

Main configuration options:

1. Postgres (started with `docker-compose -f docker-compose-dev.yml up -d`)
2. Redis (started with `docker-compose -f docker-compose-dev.yml up -d`)
3. S3 or B2
4. GitHub OAuth (https://github.com/settings/developers)
5. Google OAuth (https://console.developers.google.com/)
6. Facebook OAuth (https://developers.facebook.com/apps/)
7. Paseto keys (generated via `go run cmd/paseto/main.go`)
8. Frontend URL (needed for Google OAuth, otherwise can be ignored)
9. VirusTotal API key

The config format can be seen in `config/config.go` (each dot means a new level of nesting).