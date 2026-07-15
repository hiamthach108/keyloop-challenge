# Keyloop Inventory Backend

Go backend for Scenario B of the Keyloop technical assessment. The service is
initialized from the local Dreon backend-service boilerplate and implements
dealership inventory listing/filtering, aging-stock detection, and persisted
manager actions.

## Stack

- Go 1.25
- Echo HTTP server
- Uber Fx dependency injection
- PostgreSQL with GORM
- Dreon SDK structured logging and application errors
- Swaggo Swagger UI

Redis, gRPC, and authentication were intentionally removed because they are not
needed for the Scenario B scope.

## Run directly

From the repository root, use the root `Makefile`:

```sh
make setup
make infra-up
make backend-run
```

The API listens on `http://localhost:8080`. `GET /ping` verifies that the
service is running. Core routes are:

```text
GET  /api/v1/dealerships
GET  /api/v1/dealerships/:dealershipID/stocks
GET  /api/v1/dealerships/:dealershipID/stocks/aging
POST /api/v1/dealerships/:dealershipID/stocks/:stockID/actions
POST /api/v1/dealerships/:dealershipID/stocks/:stockID/movements
GET  /api/v1/dealerships/:dealershipID/stocks/:stockID/history
```

Swagger UI is available at:

```text
http://localhost:8080/swagger/index.html
```

The stock list supports `search` (make/model), `make`, `model`, `status`, age
filters, pagination, and enum-constrained `sortBy`/`sortOrder` values. Stock
movements and manager actions are immutable history records.

The GORM models in `internal/model/` are the schema source of truth. Atlas reads
them through `atlas.hcl` and generates versioned schema SQL with:

```sh
make migration-diff name=describe_change
```

The schema and seed migrations are stored separately in `migrations/`, and
`atlas.sum` protects both. Docker Compose runs the dedicated migration image
before starting the backend.

Regenerate Swagger after API or DTO changes:

```sh
make swagger
```

## Checks

```sh
make backend-test
make backend-lint
```
