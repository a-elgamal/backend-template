# My Service

A production-ready backend service template with a Flutter admin portal. Supports deployment to both **GCP** (Cloud Run) and **AWS** (ECS Fargate).

## Prerequisites

- Go 1.26+
- PostgreSQL 18
- Flutter 3.41+ (for portal development)
- Docker & Docker Compose

## Quick Start

```shell
# Start local Postgres
docker-compose -f docker-compose.local.yml up -d postgres

# Run all checks (lint + build + test)
make

# Run the service
make run
```

The service starts on `http://localhost:8080`. The portal is served at `http://localhost:8080/internal/portal`.

Alternatively, run everything in Docker:

```shell
docker-compose -f docker-compose.local.yml up
```

## Using This Template

### 1. Rename the module

Replace the Go module path throughout the project:

```shell
# Update go.mod
sed -i '' 's|alielgamal.com/myservice|yourorg.com/yourservice|g' go.mod

# Update all Go import paths
find . -name '*.go' -exec sed -i '' 's|alielgamal.com/myservice|yourorg.com/yourservice|g' {} +
```

### 2. Rename the service

Search and replace `myservice` with your service name in:

- `Makefile` (binary name in `run` target)
- `Dockerfile` (binary paths)
- `docker-compose.local.yml` (service name, database name)
- `application.yaml` (database name in `DB.URL_TEMPLATE`)
- `main.go` and `cmd/` (CLI command names)
- `internal/` (database table names, tracer names)
- `terraform/` (resource names, project references)
- `.github/workflows/` (image names, service references)

### 3. Replace the sample entity

The template includes a sample `app` entity in `internal/app/` demonstrating CRUD operations built on `stored.Store[T]`. To add your own domain entities:

1. Create a new package under `internal/` (e.g., `internal/order/`)
2. Define your entity struct
3. Create handlers using `stored.NewSQLStore[YourEntity](db, "table_name")`
4. Register routes in `cmd/start.go` (similar to how `app.SetupRoutes` is called)
5. Add a database migration in `internal/db/migrations/`

### 4. Enable CI/CD workflows

The template ships with CI/CD workflows that are safe to run out of the box. Test workflows (`test-server.yml`, `test-portal.yml`, `test-integration.yml`) run on every push — badge updates are skipped unless the `GIST_TOKEN` secret is configured.

Publish and deploy workflows are set to `workflow_dispatch` only (manual trigger) so they don't fail without cloud credentials. To enable them:

1. **Publish workflows** (`publish.yml`, `publish-aws.yml`): Uncomment the `push` trigger and `paths-ignore` block
2. **Deploy workflows** (`deploy-dev.yml`, `deploy-repo.yml`, `deploy-aws-dev.yml`): Uncomment the `push` and/or `workflow_run` triggers
3. Configure the required repository variables and secrets (see [GitHub Actions Setup](#github-actions-setup))

### 5. Configure deployment

Search for `TODO` comments across the Terraform and workflow files for values that need customization (project IDs, domains, VPC references, etc.).

## Project Structure

```
├── main.go                          # Entrypoint
├── cmd/                             # CLI commands (start, migrate, version)
├── internal/
│   ├── auth/                        # Auth provider interface
│   ├── aws/                         # AWS ALB OIDC auth middleware
│   ├── google/                      # GCP IAP auth middleware
│   ├── config/                      # Configuration (Viper, YAML + env vars)
│   ├── app/                         # Sample domain entity (CRUD handlers)
│   │   └── test/                    # Integration tests
│   ├── stored/                      # Generic JSONB storage layer
│   ├── db/                          # Database interfaces + migrations
│   │   ├── migrations/              # SQL migration files
│   │   └── test/                    # DB test helpers
│   ├── health/                      # Health check routes
│   ├── response/                    # Standardized error responses
│   ├── telemetry/                   # Logging (zerolog) + OpenTelemetry
│   └── testutil/                    # Integration test setup helpers
├── portal/                          # Flutter web admin portal
├── terraform/
│   ├── modules/myservice/           # GCP Terraform module
│   ├── modules/myservice-aws/       # AWS Terraform module
│   ├── dev/                         # GCP dev environment
│   ├── repo/                        # GCP Artifact Registry
│   ├── aws-dev/                     # AWS dev environment
│   └── aws-repo/                    # AWS ECR repository
├── .github/workflows/               # CI/CD pipelines
├── application.yaml                 # Default configuration
├── Dockerfile                       # Multi-stage build (Go + Flutter)
└── docker-compose.local.yml         # Local development
```

## Architecture

### Route Groups

The server organizes routes into three groups configured in `cmd/start.go`:

| Group | Path | Auth | Purpose |
|-------|------|------|---------|
| Public | `/` | None | Health checks |
| Internal | `/internal/` | GCP IAP or AWS ALB OIDC | Admin APIs + Flutter portal |
| External | `/external/` | None (add your own) | Public-facing APIs |

### Middleware Stack

Requests pass through: Telemetry metrics → OpenTelemetry tracing → Recovery → CORS → Auth (on internal routes only).

### Storage Layer

`internal/stored/` provides a generic `Stored[T]` wrapper that adds metadata (ID, timestamps, created/modified by) to any Go struct and persists it as JSONB in PostgreSQL with GIN indexes. Create a store for any type:

```go
store := stored.NewSQLStore[MyEntity](db, "my_table")
```

The store interface supports Create, Get, Update, Delete, and List with filtering.

## Commands

```shell
make                # Run lint + build + test (enforces 70% coverage minimum)
make build          # Build binary to bin/
make lint           # Run revive linter
make test           # Run tests with coverage
make fmt            # Format code with goimports
make run            # Build and run the server
```

### CLI

```shell
bin/myservice start              # Start the server
bin/myservice migrate            # Run all pending migrations
bin/myservice migrate 3          # Migrate to specific version
bin/myservice migrate --force 3  # Force version without running migration
bin/myservice version            # Print version info
```

### Running Tests

```shell
# Run all tests
go test ./...

# Run a specific test
go test ./internal/app/ -run TestHandlerName

# Run integration tests (requires local Postgres)
go test ./internal/app/test/ -run TestIntegration
```

## Configuration

Configuration is loaded from `application.yaml` and can be overridden with environment variables. Replace dots with underscores for env var names (e.g., `SERVER.HTTP_ADDRESS` → `SERVER_HTTP_ADDRESS`).

| Key | Description | Default |
|-----|-------------|---------|
| `SERVER.HTTP_ADDRESS` | Listen address | `:8080` |
| `SERVER.CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `["http://localhost:8000"]` |
| `SERVER.PORTAL_PATH` | Path to Flutter build output | `portal/build/web` |
| `SERVER.SHUTDOWN_TIMEOUT_SECONDS` | Graceful shutdown timeout | `10` |
| `DB.URL_TEMPLATE` | Postgres connection URL template | (local default) |
| `DB.NAME` | Database name | `myservice` |
| `GCP.PROJECT_NUMBER` | GCP project number | (empty) |
| `GCP.REGION` | GCP region | (empty) |
| `GCP.INTERNAL_BACKEND_SERVICE_ID` | Enables GCP IAP auth when set | (empty) |
| `AWS.ALB_REGION` | Enables AWS ALB OIDC auth when set | (empty) |
| `TELEMETRY.TRACING.SAMPLING` | Trace sampling rate (0-1) | `1` |
| `TELEMETRY.LOGGING.LEVEL` | Log level | `debug` |
| `TELEMETRY.LOGGING.CONSOLE_LOGGING_ENABLED` | Enable console logging | `TRUE` |

## Authentication

Authentication is applied only to `/internal/*` routes. The template supports two cloud providers — only one should be configured per deployment.

### GCP (Identity-Aware Proxy)

Set `GCP.INTERNAL_BACKEND_SERVICE_ID` to enable. IAP handles the OAuth flow at the load balancer level and sends a signed JWT in the `x-goog-iap-jwt-assertion` header. The middleware validates this JWT and extracts user identity.

### AWS (ALB + OIDC with Google)

Set `AWS.ALB_REGION` to enable. The ALB is configured with an OIDC authenticate action using Google as the identity provider. The ALB handles the OAuth2/OIDC flow and forwards user identity in the `x-amzn-oidc-data` header (an ALB-signed JWT). The middleware validates this JWT using ALB's public keys and extracts user identity.

In both cases, the authenticated user's ID and email are set on the Gin context and available via `internal.UserFromGinContext(c)`.

## Deployment

The same Docker image is used for both cloud providers.

### GCP (Cloud Run)

Infrastructure is defined in `terraform/modules/myservice/`:

- **Compute**: Cloud Run with an OTLP collector sidecar
- **Database**: Cloud SQL PostgreSQL 15
- **Auth**: IAP on the internal backend service
- **Registry**: Artifact Registry (configured in `terraform/repo/`)
- **Telemetry**: Exported to Google Cloud via the `googlecloud` OTLP exporter

CI/CD (manual trigger by default — see [Enable CI/CD workflows](#4-enable-cicd-workflows)):
- `publish.yml` — Builds and pushes to Artifact Registry
- `deploy.yml` — Reusable Terraform plan/apply workflow
- `deploy-dev.yml` — Deploys dev on publish completion

### AWS (ECS Fargate)

Infrastructure is defined in `terraform/modules/myservice-aws/`:

- **Compute**: ECS Fargate with an OTLP collector sidecar
- **Database**: RDS PostgreSQL 15
- **Auth**: ALB with OIDC authenticate action (Google as identity provider)
- **Registry**: ECR (configured in `terraform/aws-repo/`)
- **Telemetry**: Exported to CloudWatch/X-Ray via `awsemf`/`awsxray` OTLP exporters

CI/CD (manual trigger by default — see [Enable CI/CD workflows](#4-enable-cicd-workflows)):
- `publish-aws.yml` — Builds and pushes to ECR
- `deploy-aws.yml` — Reusable Terraform plan/apply workflow
- `deploy-aws-dev.yml` — Deploys dev on publish completion

### GitHub Actions Setup

Both cloud providers authenticate via GitHub OIDC (no long-lived secrets):

**GCP** — Set these repository variables:
- `GCP_PROJECT_ID`, `GCP_WORKLOAD_IDENTITY_PROVIDER`, `GCP_SA`
- `TERRAFORM_VERSION`

**AWS** — Set these repository variables:
- `AWS_DEPLOY_ROLE_ARN`, `AWS_REGION`
- `TERRAFORM_VERSION`

**Badges** (optional) — Set the `GIST_TOKEN` secret to a personal access token with `gist` scope to enable test/coverage badge updates.

## Portal Development

The Flutter portal is located in `portal/` and is served at `/internal/portal` in production.

For local development:

1. Run Postgres: `docker-compose -f docker-compose.local.yml up -d postgres`
2. Run the backend: `make run`
3. Run the portal with `flutter run` from the `portal/` directory
4. Add the Flutter dev server URL to `SERVER.CORS_ALLOWED_ORIGINS` in `application.yaml`

## Database Migrations

Migrations use [golang-migrate](https://github.com/golang-migrate/migrate) with SQL files in `internal/db/migrations/`. Files follow the naming convention `{version}_{description}.up.sql` and `{version}_{description}.down.sql`.

To add a new migration, create the next numbered pair of files in that directory and they will run automatically on service startup.
