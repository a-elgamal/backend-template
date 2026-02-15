# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go backend service template with a Flutter web portal. It uses the `stored` generic JSONB layer for data persistence and includes a sample `item` domain entity as a working example.

## Common Commands

```shell
# Prerequisites: start local Postgres
docker-compose -f docker-compose.local.yml up

# Run all checks (lint + build + test)
make

# Individual targets
make build          # Build binary to bin/
make lint           # Run revive linter
make test           # Run tests with coverage (enforces 70% minimum)
make fmt            # Format code with goimports
make run            # Build and run the server

# Run a single test
go test ./internal/item/ -run TestHandlerName

# Run integration tests (requires local Postgres)
go test ./internal/item/test/ -run TestIntegration

# Database migrations
bin/myservice migrate          # Run all pending migrations
bin/myservice migrate 3        # Migrate to specific version
bin/myservice migrate --force 3  # Force version without running migration
```

## Architecture

### Backend (Go + Gin)

**Entrypoint**: `main.go` -> `cmd/` (Cobra CLI with `start`, `migrate`, `version` subcommands)

**Route groups** (configured in `cmd/start.go`):
- `/` — Public health check
- `/internal/` — Admin APIs + Flutter portal (Google IAP auth when GCP configured)
- `/external/` — External APIs

**Middleware stack**: Telemetry metrics -> OpenTelemetry tracing (otelgin) -> Recovery -> CORS -> Auth (IAP)

### Key Packages

- **`internal/stored/`** — Generic JSON storage layer. `Stored[T]` wraps any type with metadata (ID, timestamps, created/modified by). `sqlStore[T]` persists content as JSONB in Postgres with GIN indexes.
- **`internal/item/`** — Sample domain entity (CRUD handlers). Built on top of `stored.Store[Item]`.
- **`internal/db/`** — Database interfaces (`DB`, `Tx`) wrapping `sql.DB`/`sql.Tx`. Migrations use `golang-migrate` with embedded SQL files in `internal/db/migrations/`.
- **`internal/config/`** — Config loaded from `application.yaml` via Viper with environment variable overrides.
- **`internal/telemetry/`** — Zerolog logging + OpenTelemetry traces/metrics setup.
- **`internal/google/`** — IAP JWT validation middleware.
- **`internal/response/`** — Standardized error response structs.

### Frontend (Flutter/Dart)

Located in `portal/`. Uses GoRouter for routing, Provider for state management. Served at `/internal/portal`.

## Testing Patterns

- **Unit tests**: Mock the `Store[T]` interface using testify mocks. Test handlers with `httptest.NewRecorder()`.
- **Integration tests**: Located in `test/` subdirectories. Use `testutil.SetUpIntegartionTest()` which starts a real server with an isolated test database.
- **DB test helpers**: `internal/db/test/` provides `SetupTestDB()` for creating isolated test databases.

## Configuration

`application.yaml` at the project root. Environment variables override YAML values (dots become underscores).
