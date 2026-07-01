# Go Template

A clean-architecture template for a Go REST microservice. The Go module is named `application`.

## What's in the box

- **HTTP server** — stdlib `net/http.ServeMux` with `otelhttp` + panic-recovery middleware (`app/httpserver.go`).
- **Compile-time DI** — Google Wire; composition root in `cmd/app/wire.go`.
- **Configuration** — koanf, YAML file plus `APP_`-prefixed env overlay (`app/config.go`).
- **Observability** — OpenTelemetry traces, metrics, and logs over OTLP gRPC; slog bridged to OTel logs (`app/otel.go`, `app/logger.go`).
- **Datasources** — PostgreSQL (pgx + otelsql), an in-memory ramsql DB for testing, and a NATS/JetStream client.
- **Lifecycle** — components self-register `Start`, `Shutdown`, and `Healthz` hooks on a shared `app.Controller` (`app/controller.go`).
- **Mocks** — mockery v2 generates from `internal/biz` into `internal/mocks` (`.mockery.yaml`).
- **API docs** — swag generates Swagger; UI mounted at `/swagger/`.

## Quickstart

```sh
cp config.example.yaml config.yaml
make devtools                      # one-time: golangci-lint, gofumpt, mockgen, swag, gci
make generate                      # tidy go.mod, install wire, run go generate
go run ./cmd/app --config ./config.yaml
```

The server listens on `:8080` by default. Override any config key via env, e.g. `APP_SERVER_HTTP_ADDR=:9090`.

## Endpoints exposed by default

| Path | Notes |
|---|---|
| `GET /healthz/liveness` | Liveness check (fans out across registered checks) |
| `GET /healthz/readiness` | Readiness check |
| `GET /healthz/panic` | Triggers a panic for testing the recovery middleware |
| `GET /healthz/sleep/{time}` | Sleeps for a `time.Duration` (e.g. `10s`) |
| `GET /apis/mocks/placeholders` | Example CRUD resource (list/get/create/update/delete) |
| `GET /metrics` | Prometheus exposition |
| `GET /swagger/` | Swagger UI (spec at `/docs/swagger/swagger.json`) |

## Common commands

```sh
make generate         # regenerate Wire DI graph + go.mod tidy
make swagger          # regenerate ./docs/ from swag annotations
make check            # golangci-lint
make unit_tests       # tests under internal/service/handler and internal/biz
make coverage_tests   # same, with coverage profile to coverage.out
make all_tests        # tests + benchmarks + coverage (JSON output)
make build            # docker build -t buildf .
```

Run a single test:

```sh
go test ./internal/service/handler/... -run TestName -v
```

## Architecture

```
cmd/app                 main + wire composition root
app/                    runtime infra (Application, HTTPServer, Controller, KConfig, Logger, OTLP)
internal/service        HTTP wiring (mux, /metrics, /swagger)
internal/service/handler HTTP handlers — implement service.Handler
internal/service/dto    request/response shapes + error mapping
internal/biz            use cases (interfaces consumed by handlers)
internal/repo           repository implementations (bound to biz Repository interfaces)
internal/datasource     DB / queue clients
internal/entity         domain types
internal/mocks          mockery-generated mocks for internal/biz
pkg/middlewares         per-route HTTP middlewares
infra/                  docker-compose stacks (monitoring, postgres, redis) + k8s manifests
migrations/             SQL migrations
```

Dependency direction: `cmd → app → service/handler → biz → repo → datasource → entity`.

To add a handler: implement `service.Handler` (`RegisterHandler(ctx) error` registers routes on the injected `*http.ServeMux`), expose a `New…` provider, and append it to `NewServiceList` in `internal/service/handler/wire.go`. Then run `make generate`.

## Configuration

`config.yaml` is the source of truth at runtime; copy from `config.example.yaml`. Env vars prefixed `APP_` are merged on top, with `_` mapping to `.` (e.g. `APP_SERVER_HTTP_ADDR` → `server.http.addr`). The `--config` flag selects the file (default `./config.yaml`).

## Local infrastructure

`docker-compose.yml` aggregates includes from `infra/compose/` for monitoring, Postgres, Redis, and MinIO:

```sh
docker compose up -d
```

The `Tiltfile` references a Helm chart at `./charts/...` that is not included in this repo — add your own chart to use Tilt.

---

## BOTB backend (monorepo)

This template is being grown into a **single-module monorepo** (`module application`)
that serves both the public competition site and the admin panel from one domain.
Each service is its own deployable binary under `cmd/`, all reusing the shared
`app/` + `pkg/` infrastructure. Routes follow `/apis/<service>/<version>/...`.
Public reads are unauthenticated; admin mutations live under an `/admin/` route
group guarded by JWT at the gateway.

| Service | Binary | Status |
|---|---|---|
| Media (uploads, object storage) | `cmd/media` | ✅ implemented |
| Competition (public GET + admin CRUD) | `cmd/competition` | ✅ implemented |
| User + Tickets | `cmd/user` | ✅ implemented |
| Draw / Winners | `cmd/draw` | ✅ implemented |
| Gateway (single public entrypoint, reverse-proxy + JWT) | `cmd/gateway` | ✅ implemented |

### Media service (`cmd/media`)

Real file uploads (image or video) to S3-compatible object storage (MinIO now,
AWS S3 later via the same `ObjectStorage` port). Metadata lives in Postgres so
any object — e.g. a competition — can have zero, one, or many media items,
resolved by `owner_type` + `owner_id`.

| Method & Path | Auth | Description |
|---|---|---|
| `POST /apis/media/v1/uploads` | admin | Upload a file (multipart: `file`, `owner_type`, `owner_id`, `position`) |
| `GET /apis/media/v1/media/{id}` | public | Media metadata + a presigned read URL |
| `GET /apis/media/v1/media?owner_type=&owner_id=` | public | List all media for an owner |

Validation: images `jpeg/png/webp` ≤ 10 MB, videos `mp4/webm` ≤ 200 MB.

### Competition service (`cmd/competition`)

Public reads + admin CRUD. Each competition returns its associated media,
resolved by a read query against the shared `media` table (owner_type +
owner_id) — the most template-consistent choice given the single pgx datasource.

| Method & Path | Auth | Description |
|---|---|---|
| `GET /apis/competition/v1/competitions` | public | List competitions (`?status=draft\|live\|closed`) |
| `GET /apis/competition/v1/competitions/{id}` | public | One competition + its media |
| `POST /apis/competition/v1/admin/competitions` | admin | Create |
| `PUT /apis/competition/v1/admin/competitions/{id}` | admin | Update |
| `DELETE /apis/competition/v1/admin/competitions/{id}` | admin | Delete |

Admin routes live under the `/admin/` segment so the gateway can guard that
group with JWT. Money is stored as integer pence (`ticket_price_pence`).

### User + Ticket service (`cmd/user`)

Owns the `users` and `tickets` tables. Public registration + ticket purchase;
admin user listing/lookup.

| Method & Path | Auth | Description |
|---|---|---|
| `POST /apis/user/v1/users` | public | Register (name + email) |
| `POST /apis/user/v1/tickets` | public | Purchase `quantity` tickets for a competition |
| `GET /apis/user/v1/admin/users` | admin | Searchable, paginated user list (`?q=&limit=&offset=`) |
| `GET /apis/user/v1/admin/users/{id}` | admin | One user |
| `GET /apis/user/v1/admin/users/{id}/tickets` | admin | A user's tickets |

Purchase runs in a transaction: it reads the competition's `ticket_price_pence`
from the shared DB, inserts the ticket rows, and atomically increments the
user's `tickets_owned` + `total_spent_pence`. It does **not** write
`competitions.tickets_sold` (owned by the competition service — a real system
would sync that via a JetStream event). Registration (a public `POST`) is the
one deliberate exception to "writes are admin-only".

### Draw / Winners service (`cmd/draw`)

Owns the `draws` table. Admin creates + runs draws; the public homepage reads
completed results.

| Method & Path | Auth | Description |
|---|---|---|
| `GET /apis/draw/v1/admin/draws` | admin | Searchable, paginated draw list (`?q=&limit=&offset=`) |
| `GET /apis/draw/v1/admin/draws/{id}` | admin | One draw (incl. pending) |
| `POST /apis/draw/v1/admin/draws` | admin | Create a pending draw for a competition |
| `POST /apis/draw/v1/admin/draws/{id}/run` | admin | Run the draw — pick a winner |
| `GET /apis/draw/v1/draws/{id}` | public | Read a completed result (pending → 404) |

Running a draw is a single transaction: it guards that the draw is still
`pending` (inside the tx, with a conditional `WHERE ... AND status = 'pending'`
UPDATE so concurrent runs can't both win), reads the competition's tickets from
the shared DB, and picks a uniformly-random winner with **crypto/rand** (not
math/rand). Like the ticket-purchase flow it does not mutate competition-owned
state (a JetStream event would sync competition status in a real system).

### Gateway service (`cmd/gateway`)

The single public entrypoint. It holds no state — it reverse-proxies (stdlib
`httputil.NewSingleHostReverseProxy`) by the `<servicename>` path segment of
`/apis/<servicename>/...` to that service's upstream base URL (from
`gateway.upstreams.*`, `APP_`-overridable). Unknown service → `404`. Trace
context is propagated to upstreams via an otel-instrumented transport.

**Two-layer (defense-in-depth) admin auth.** A shared HS256 bearer-token
middleware (`pkg/middlewares/jwtauth.go`) guards any path matching
`/apis/<svc>/v1/admin/...`:

1. at the **gateway**, before proxying, and
2. **inside each service** (its own admin route group), so a service reached
   directly on its internal port is never unprotected.

Both read the **same** secret from `jwt.secret` (one definition, shared infra
config). Public reads and the public POSTs (register, purchase) need no token;
missing/invalid/expired token on an admin path → `401`.

### Running the full stack locally

```sh
docker compose up -d postgresql minio          # infra
cp config.example.yaml config.yaml
# apply migrations/ against Postgres (e.g. with golang-migrate)

# each domain service on its own port (matches gateway.upstreams in config)
APP_SERVER_HTTP_ADDR=:8081 go run ./cmd/competition --config ./config.yaml &
APP_SERVER_HTTP_ADDR=:8082 go run ./cmd/user        --config ./config.yaml &
APP_SERVER_HTTP_ADDR=:8083 go run ./cmd/draw        --config ./config.yaml &
APP_SERVER_HTTP_ADDR=:8084 go run ./cmd/media       --config ./config.yaml &

# gateway on the single public port
APP_SERVER_HTTP_ADDR=:8080 go run ./cmd/gateway --config ./config.yaml
```

All public traffic then goes through `http://localhost:8080/apis/...`. Service
binaries run via `go run` (not in docker-compose, which stays infra-only).

### New config / env vars

`config.example.yaml` gained a config-driven Postgres DSN, a MinIO block, the
shared `jwt.secret`, and `gateway.upstreams`. All keys are overridable via the
`APP_` env convention (`_` → `.`):

| Config key | Env override | Default |
|---|---|---|
| `datasource.postgres.dsn` | `APP_DATASOURCE_POSTGRES_DSN` | `postgresql://botb:botb@localhost:5432/botb?sslmode=disable` |
| `datasource.minio.endpoint` | `APP_DATASOURCE_MINIO_ENDPOINT` | `localhost:9000` |
| `datasource.minio.access_key` | `APP_DATASOURCE_MINIO_ACCESS_KEY` | `minioadmin` |
| `datasource.minio.secret_key` | `APP_DATASOURCE_MINIO_SECRET_KEY` | `minioadmin` |
| `datasource.minio.bucket` | `APP_DATASOURCE_MINIO_BUCKET` | `botb-media` |
| `datasource.minio.use_ssl` | `APP_DATASOURCE_MINIO_USE_SSL` | `false` |
| `datasource.minio.presign_expiry` | `APP_DATASOURCE_MINIO_PRESIGN_EXPIRY` | `15m` |
| `jwt.secret` | `APP_JWT_SECRET` | `dev-insecure-change-me` |
| `gateway.upstreams.<svc>` | `APP_GATEWAY_UPSTREAMS_<SVC>` | `http://localhost:808x` |

> Note: the `make check` target references a `v1`-style flag set and a missing
> `issues.exclude.yaml`; per `CLAUDE.md`, run `golangci-lint run` directly (it
> uses `.golangci.yaml`). The whole module currently passes with 0 issues.
