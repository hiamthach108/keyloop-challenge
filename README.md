# Keyloop Intelligent Inventory Dashboard

Scenario B implementation for the Keyloop technical assessment. The product
gives dealership managers a filterable inventory view, identifies vehicles held
for more than 90 days, and persists proposed actions for aging stock.

## Repository layout

```text
.
├── ARCHITECTURE.md
├── Makefile
├── backend/         Go REST API and PostgreSQL persistence
├── docker-compose.yml
└── frontend/        Next.js and Ant Design dashboard
```

The frontend uses the real dealership-scoped backend API by default. A mock
adapter remains available for isolated UI development. Authentication is
intentionally outside this challenge's scope.

## Prerequisites

- Go 1.25 or newer
- Docker with Docker Compose
- `make`

## Setup and run

The simplest full-stack command builds and runs PostgreSQL, the backend, and
the frontend in Docker:

```sh
make docker-up
```

Open `http://localhost:3000`. The API is available at
`http://localhost:8080/api/v1`, and Swagger UI is available at
`http://localhost:8080/swagger/index.html`. Compose runs a dedicated Atlas
migration container before starting the backend. Atlas applies the generated
schema migration and the checksum-protected seed migration under
`backend/migrations/`.

Run the stack in the background with `make docker-up-detached`, and stop it
with `make docker-down`. To recreate and reseed the database:

```sh
make docker-reset
```

For native development with only PostgreSQL in Docker:

```sh
make setup
make infra-up
make backend-run
```

The backend listens on `http://localhost:8080`. Verify it with:

```sh
curl http://localhost:8080/ping
```

To start PostgreSQL, the backend, and the frontend together:

```sh
make dev
```

The dashboard is available at `http://localhost:3000` and uses the backend API.
Set `NEXT_PUBLIC_INVENTORY_SOURCE=mock` in `frontend/.env.local` only when you
want to run the UI independently.

Stop PostgreSQL with:

```sh
make infra-down
```

## Development checks

```sh
make backend-test
make backend-lint
make backend-tidy
make frontend-lint
make frontend-typecheck
make frontend-build
```

## Database migrations

GORM models under `backend/internal/model/` are the schema source of truth.
Generate a migration after changing a model with:

```sh
make migration-diff name=describe_change
```

Seed SQL is maintained as a separate versioned migration. After changing it,
refresh and verify the Atlas checksum:

```sh
make migration-hash
make migration-validate
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for the system design, assumptions, API
shape, data flow, observability strategy, and delivery plan.

## Swagger

Swagger annotations live beside the Echo handlers and generated docs are stored
under `backend/docs/`.

Regenerate the spec after changing backend routes or request/response DTOs:

```sh
make swagger
```

Then open `http://localhost:8080/swagger/index.html` while the backend is
running.

## AI collaboration narrative

GenAI is being used as a directed implementation partner. Requirements and
ambiguities are first translated into explicit architecture and acceptance
criteria. Generated changes are then reviewed against dealership isolation,
the 90-day business boundary, persistence behavior, API contracts, and tests.
The final implementation remains developer-owned: dependencies are kept
purposeful, unused boilerplate is removed, and every core behavior will be
verified through automated tests and an end-to-end API walkthrough.
