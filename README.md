# SMR API [![build](https://github.com/satisfactorymodding/smr-api/actions/workflows/build.yml/badge.svg)](https://github.com/satisfactorymodding/smr-api/actions/workflows/build.yml) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/satisfactorymodding/smr-api) [![codecov](https://codecov.io/gh/satisfactorymodding/smr-api/branch/master/graph/badge.svg?token=LFNKYWS0N2)](https://codecov.io/gh/satisfactorymodding/smr-api) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/satisfactorymodding/smr-api)

The Satisfactory Mod Repository backend API - a Go-based service providing REST and GraphQL endpoints for mod management.

## Quick Start

### Prerequisites

- [mise](https://mise.jdx.dev/) for tool management
- [Docker (and Docker Compose)](https://docs.docker.com/desktop/)

### Setup

1. Install Tools

   ```bash
   # Install tools
   mise install
   ```

   If you are on Windows, due to a bug in mise, you will also need to do the following:

   ```ps1
   > mise where aqua:minio/mc
   <some folder path>
   # Go to the folder
   # Rename the `mc` file to give it an extension: `mc.exe`
   ```

2. Start Development Services

   ```bash
   # If this command fails, ensure you have Docker running and configured correctly
   mise run setup

   # Generate code and start API
   mise run generate
   mise run api
   ```

   If running `api` produces errors about the database failing to apply migrations,
   you may have switched branches without cleaning up after database changes you were working on.
   The easiest way to get back from this state
   is to delete the `postgres` Docker container
   and delete all not-in-use volumes so it creates a fresh one.

3. Set up the Configuration File

   Create a copy of `config.sample.json` as `config.json` and fill in the required fields.
   See the [Configuration](#configuration) section for details.

## Development Commands

```bash
# Code generation (run after schema changes)
mise run generate

# Start API server
mise run api

# Testing
mise run localtest
mise run coverage

# Linting
mise run lint

# Environment management
mise run setup     # Start PostgreSQL, Redis, MinIO
mise run teardown  # Stop services

# Windows: activate an interactive PowerShell with Mise tools loaded. Requires Powershell >= 7
mise activate pwsh | Out-String | Invoke-Expression
```

### Local Testing

To run specific tests:

- In VSCode, the [Go Companion](https://marketplace.visualstudio.com/items?itemName=ethan-reesor.exp-vscode-go) suggested extension adds Go test support to VSCode's testing integration.
  Note you may need to reload VSCode (`Developer: Reload Window`) after adding or removing tests.
- From the command line, use `go test -v run TestNameHere`

To run all tests on your local machine, use `mise run localtest`, which excludes the verbose flag for easier human parsing of test output.
It is expected for the command to output `[no test files]` for packages with no tests.

### Database Access

The development docker compose also includes a pgadmin container for testing database related stuff
exposed at <http://localhost:5433/>.
The credentials are specified in `docker-compose-dev.yml`.

Once logged in, add a server.
The "Host name/address" should be `postgres`, username should be `postgres`, password should be `POSTGRES_PASSWORD` from `docker-compose.yml`.

You can view tables via `Servers >(name)> Databases > postgres > Schemas > public > Tables`.

## Architecture

**Tech Stack:**

- **Framework**: Echo v4 with middleware
- **Database**: PostgreSQL with Ent ORM
- **GraphQL**: gqlgen with custom resolvers
- **Authentication**: Multi-provider OAuth + PASETO tokens
- **Storage**: S3-compatible (MinIO for dev)
- **Background Processing**: Temporal.io workflows
- **Monitoring**: OpenTelemetry tracing

**Key Components:**

- REST API at `/v1/`
- GraphQL at `/v2/` with playground
- Swagger docs at `/swagger/`
- Generated code in `generated/`
- Database schemas in `db/schema/`

## Configuration

Create `config.json`, or use environment variables with `REPO_` prefix.
A sample is provided: `config.sample.json`.

**Required services:**

1. **PostgreSQL** - Database (dev: port 5432)
2. **Redis** - Cache/sessions (dev: port 6379)  
3. **MinIO** - Object storage (dev: ports 9000/9001)
4. **OAuth providers** - GitHub, Google, Facebook. For testing purposes, [setting up just GitHub](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app) is the easiest.
5. **PASETO keys** - Generate these with `go run cmd/paseto/main.go`
6. **VirusTotal API key** - For mod scanning. Optional for testing.

**Development services** are started automatically with `mise run setup`. MinIO configuration is included in the setup task.

See `config/config.go` for full configuration structure.

## Development Workflow

1. **Code Generation**: Always run `mise run generate` after modifying:
   - GraphQL schemas (`schemas/*.graphql`)
   - Database schemas (`db/schema/*.go`)
   - Swagger annotations

2. **Database Changes**: Use Atlas migrations
   - SQL migrations in `migrations/sql/`
   - Code migrations in `migrations/code/`

3. **Testing**: Tests require development services running

   ```bash
   mise run setup  # Start services
   mise run test   # Run tests
   ```

## Contributing

**Before submitting:**

```bash
mise run format   # Automatically fix code quality problems and report those that need manual correction
mise run test     # Run test suite
mise run generate # Regenerate if needed
```

**Development patterns:**

- Use Ent ORM for database operations
- Implement GraphQL resolvers for complex queries
- Use Temporal workflows for background processing
- Follow structured logging with `slog`
- Use context for request tracing
