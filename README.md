# SMR API [![build](https://github.com/satisfactorymodding/smr-api/actions/workflows/build.yml/badge.svg)](https://github.com/satisfactorymodding/smr-api/actions/workflows/build.yml) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/satisfactorymodding/smr-api) [![codecov](https://codecov.io/gh/satisfactorymodding/smr-api/branch/master/graph/badge.svg?token=LFNKYWS0N2)](https://codecov.io/gh/satisfactorymodding/smr-api) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/satisfactorymodding/smr-api)

The Satisfactory Mod Repository backend API - a Go-based service providing REST and GraphQL endpoints for mod management.

## Quick Start

### Prerequisites

- [mise](https://mise.jdx.dev/) for tool management
- [Docker (and Docker Compose)](https://docs.docker.com/desktop/)
- If on Windows, the [latest PowerShell (>= 7)](https://learn.microsoft.com/en-us/powershell/scripting/install/install-powershell-on-windows). Check current version by running `pwsh -v` in a PowerShell window. ([reason why](https://github.com/jdx/mise/discussions/4151))

### Setup

> Note: A [devcontainer](https://code.visualstudio.com/docs/devcontainers/containers) exists in this repository,
> but it has not been tested since we switched to using Mise for tool management.
> **You should follow the normal installation instructions below instead of using the devcontainer.**
> We offer no guarantee that it works correctly, but you are welcome to try it out and submit fixes if needed.

1. Install Tools

   ```bash
   # Install tools
   mise install
   ```

   If you are on Windows, due to a bug in [mise and aqua](https://github.com/aquaproj/aqua-registry/pull/42498/), you will also need to do the following:

   <!-- TODO make a mise task that does this for you -->
   ```ps1
   > mise where aqua:minio/mc
   <some folder path>
   # Go to the folder
   # Rename the `mc` file to give it an extension: `mc.exe`
   ```

2. (First Time Setup Only) Set up the Configuration File

   First, start the containers so you can generate some of the required configuration tokens:

   ```bash
   # If this command fails, ensure you have Docker running and configured correctly
   mise run setup
   ```

   Next, create a copy of `config.sample.json` as `config.json` and fill in the required fields.
   See the [Configuration](#configuration) section for details.

3. Start Development Services

   ```bash
   # If this command fails, ensure you have Docker running and configured correctly
   mise run setup

   # Generate code and start API
   mise run generate
   mise run api
   ```

   If the first line of the output contains `config initialized using defaults and environment only!` then ensure your config file is formatted correctly and does not include `//` comments.

   If running `api` produces errors about the database failing to apply migrations,
   you may have switched branches without cleaning up after database changes you were working on.
   The easiest way to get back from this state
   is to delete the `postgres` Docker container
   and delete all not-in-use volumes so it creates a fresh one.

## Development Commands

```bash
# Get information about all available mise tasks
mise tasks

# Environment management
mise run setup     # Start PostgreSQL, Redis, MinIO, etc. containers
mise run teardown  # Stop services

# Code generation (run after schema changes)
mise run generate

# Start API server
mise run api

# Testing
mise run localtest
mise run coverage

# Linting
mise run lint

# Linux (bash): activate an interactive bash shell with Mise tools loaded
eval "$(mise activate bash)"

# Windows: activate an interactive PowerShell with Mise tools loaded. Requires Powershell >= 7
mise activate pwsh | Out-String | Invoke-Expression

# Database migrations. See `mise tasks` for more info.
mise run migrate_diff
```

### Local Testing

To run specific tests:

- Make sure you have started the containers first via `mise run setup`
- In VSCode, the [Go Companion](https://marketplace.visualstudio.com/items?itemName=ethan-reesor.exp-vscode-go) suggested extension adds Go test support to VSCode's testing integration.
  Note you may need to reload VSCode (`Developer: Reload Window`) after adding or removing tests.
- From the command line, use `go test -v run TestNameHere`

To run all tests on your local machine, use `mise run localtest`, which excludes the verbose flag for easier human parsing of test output.
It is expected for the command to output `[no test files]` for packages with no tests.

### Database Access

The development docker compose also includes a pgadmin container for testing database related stuff
exposed at <http://localhost:5433/>.
The credentials are specified in [`docker-compose-dev.yml`](docker-compose-dev.yml).

Once logged in, add a server.
The "Host name/address" should be `postgres`, username should be `postgres`, password should be `POSTGRES_PASSWORD` from [`docker-compose.yml`](docker-compose.yml).

You can view tables via `Servers > (name) > Databases > postgres > Schemas > public > Tables`.

#### Site User Roles

To locally test features of the site that require a user role, you can change your role via direct database queries.
Doing this requires manually creating a `user_groups` entry for you user.

1. Create an account on your local via the [frontend running with the `development` env](https://github.com/satisfactorymodding/smr-frontend/blob/staging/CONTRIBUTING.md#decide-which-environment-you-want-to-run)
2. Look at [`permissions.go`](/auth/permissions.go) to find the group ID for the group you want to become. Admin is `1`.
3. Run a query in pgadmin to determine your user ID: `SELECT * FROM users;`
4. Run a query to create a user_groups entry for yourself. Example for admin:

   ```sql
   INSERT INTO public.user_groups(
   user_id, group_id, created_at, updated_at, deleted_at, id)
   VALUES ('YOUR_USER_ID_HERE', '1', NOW(), NOW(), null, NOW());
   ```

5. Refresh the frontend webpage - you should be an admin now.

## Architecture

**Tech Stack:**

- **Framework**: Echo v4 with middleware
- **Database**: PostgreSQL with Ent ORM
  - [Atlas versioned migrations](https://atlasgo.io/versioned/intro)
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

**Required services:**

1. **PostgreSQL** - Database (dev: port 5432)
2. **Redis** - Cache/sessions (dev: port 6379)  
3. **MinIO** - Object storage (dev: ports 9000/9001)
4. **VirusTotal** - External API for mod scanning. Optional for testing and local development.

## Configuration

**Development services** are started automatically with `mise run setup`. MinIO configuration is included in the setup task.

Create `config.json`, or use environment variables with `REPO_` prefix.
A sample is provided: `config.sample.json`.

**For local development, you should change:**

1. **OAuth providers** - For testing purposes, [setting up just GitHub](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app) is the easiest.
2. **PASETO keys** - Generate these with `mise run make_paseto_keys`
3. **VirusTotal API key** - If you want to test the mod scanning features (optional).

See `config/config.go` for full configuration structure.

## Development Workflow

1. **Code Generation**: Always run `mise run generate` after modifying:
   - GraphQL schemas (`schemas/*.graphql`)
   - Database schemas (`db/schema/*.go`)
   - Swagger annotations

2. **Database Changes**: Use Atlas migrations (entgo has integration with it via golang-migrate)
   - SQL migrations in `migrations/sql/`
     - Use mise tasks to create migration files (`mise tasks` for more details)
     - Migration names should be in `snake_case`
   - Code migrations in `migrations/code/`
     - Created manually
     - Get applied after all SQL migrations (see [migrations.go](migrations/migrations.go) for details)

3. **Testing**: Tests require development services running

   ```bash
   mise run setup  # Start services

   mise run localtest   # Run tests
   # OR
   mise run test   # Run tests (verbose)
   ```

## Contributing

**Before submitting:**

Commit messages should follow the [Conventional Commits format](https://www.conventionalcommits.org/en/v1.0.0/).

```bash
mise run lint       # Check code quality
mise run format     # Some linting issues can be auto-fixed
mise run localtest  # Run test suite
mise run generate   # Regenerate if needed
```

**Development patterns:**

- Use Ent ORM for database operations
- Implement GraphQL resolvers for complex queries
- Use Temporal workflows for background processing
- Follow structured logging with `slog`
- Use context for request tracing
