# FlyPro Travel Expense Management API (Go + Gin + GORM + Postgres + Redis)

A production-style backend showing clean architecture, repository/service pattern,
Goose SQL migrations, validation, caching, and a currency conversion integration.

## Stack
- Go 1.23.0
- Gin (HTTP)
- GORM (ORM) with PostgreSQL
- Redis (caching)
- Goose (SQL migrations)
- Docker Compose (dev env)

## Quickstart

```bash
cp .env.example .env
docker compose up -d --build
# run migrations
make migrate-up
# seed demo data
make seed
```

App will be available on `http://localhost:8080`.

### Make Targets
- `make run` - run app locally
- `make test` - run tests
- `make migrate-up` / `make migrate-down` / `make migrate-status` - Goose migrations
- `make new-migration name=descriptive_name` - create new SQL migration (timestamped)
- `make seed` - seed demo data

### Postman Collection
Collection at `docs/FlyPro Expense API.postman_collection.yaml`. You can import it into Postman.

### Architecture
- `cmd/server/main.go` bootstraps the app
- `internal/config` configuration
- `internal/models` GORM models (no automigrate in prod; use Goose)
- `internal/dto` request/response DTOs with validation tags
- `internal/validators` custom validators (currency code, category)
- `internal/repository` data access
- `internal/services` business logic (currency conversion + caching)
- `internal/handlers` HTTP handlers (Gin)
- `internal/middleware` logging, request ID, rate limiting, recovery
- `internal/utils` helpers (errors, pagination, cache)
- `migrations` Goose SQL migrations

### Design Notes
- **Migrations**: strictly via Goose SQL files (no GORM automigrate in server).
- **Validation**: go-playground validator via Gin binding + custom funcs.
- **Caching**: Redis for FX rates (6h TTL), users (1h), report summaries (30m),
  and list endpoints (5m demo TTL). Cache keys are namespaced.
- **Currency**: reports compute totals in USD using the FX service with graceful
  degradation (fallback to last cached rate).
- **Rate limiting**: per-IP token bucket using in-memory map with cleanup.
