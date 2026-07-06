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
group that requires a **valid JWT carrying an admin `role` claim** (admin or
superadmin), enforced at the gateway AND inside each service.

| Service | Binary | Status |
|---|---|---|
| Media (uploads, object storage, admin lifecycle) | `cmd/media` | ✅ implemented |
| Competition (public GET + admin CRUD, categories, site content) | `cmd/competition` | ✅ implemented |
| User + Tickets (admin edit + suspend) | `cmd/user` | ✅ implemented |
| Draw / Winners (edit, void, reassign, public winners feed) | `cmd/draw` | ✅ implemented |
| Adminauth (independent admin identity, login + accounts) | `cmd/adminauth` | ✅ implemented |
| Gateway (single public entrypoint, reverse-proxy + role guard) | `cmd/gateway` | ✅ implemented |

### Adminauth service (`cmd/adminauth`)

The independent admin identity provider. Admin accounts live in their **own
`adminauth` Postgres schema** (`adminauth.admin_accounts`,
`adminauth.admin_refresh_tokens`) with **zero foreign keys** to the site's
public schema — admins are NOT site users. Passwords are bcrypt-hashed; login
returns a generic error for every failure mode (unknown email, wrong password,
disabled account) and is rate-limited (5/min per email and per IP). It signs
short-lived HS256 access tokens with the shared `jwt.secret` (so every guard
accepts them) plus a rotating refresh token (reuse of a rotated token revokes
the whole chain).

| Method & Path | Auth | Description |
|---|---|---|
| `POST /apis/adminauth/v1/login` | public | Verify credentials → access JWT + refresh token |
| `POST /apis/adminauth/v1/refresh` | public | Rotate a refresh token into a new pair |
| `POST /apis/adminauth/v1/logout` | public | Revoke a refresh token |
| `GET /apis/adminauth/v1/me` | admin | Current admin profile + role |
| `GET/POST /apis/adminauth/v1/admin/accounts` | superadmin | List / create admin accounts |
| `PUT /apis/adminauth/v1/admin/accounts/{id}` | superadmin | Edit / disable (guards the last active superadmin) |

**First superadmin bootstrap.** On startup, if `APP_ADMINAUTH_BOOTSTRAP_EMAIL`
and `APP_ADMINAUTH_BOOTSTRAP_PASSWORD` are set AND no superadmin exists yet,
adminauth creates one (role=superadmin). Idempotent, never overwrites, never
logs the password. Compose wires these from `ADMIN_BOOTSTRAP_*` in `.env`
(defaults `admin@botb.local` / `change-me-now` — change and rotate after first
login). To bootstrap manually instead, leave those env vars blank and create
the account from a one-shot run of the service with them set.

### Media service (`cmd/media`)

Real file uploads (image or video) to S3-compatible object storage (MinIO now,
AWS S3 later via the same `ObjectStorage` port). Metadata lives in Postgres so
any object — e.g. a competition — can have zero, one, or many media items,
resolved by `owner_type` + `owner_id`.

| Method & Path | Auth | Description |
|---|---|---|
| `GET /apis/media/v1/media/{id}` | public | Media metadata + a presigned read URL |
| `GET /apis/media/v1/media?owner_type=&owner_id=` | public | List all media for an owner |
| `POST /apis/media/v1/admin/uploads` | admin | Upload a file (multipart: `file`, `owner_type`, `owner_id`, `position`) |
| `GET /apis/media/v1/admin/media` | admin | By-owner, or a paged global media library |
| `PUT /apis/media/v1/admin/media/{id}` | admin | Reorder (`position`) and/or reassign owner |
| `DELETE /apis/media/v1/admin/media/{id}` | admin | Remove DB row + MinIO object |

Validation: images `jpeg/png/webp` ≤ 10 MB, videos `mp4/webm` ≤ 200 MB. The
old **unauthenticated** `POST /uploads` is gone — uploads are admin-only now.
Delete removes the DB row first, then best-effort removes the object (a failed
object removal is a logged storage leak, not a dangling DB row). Replace =
upload new + delete old.

### Competition service (`cmd/competition`)

Public reads + admin CRUD. Each competition returns its associated media,
resolved by a read query against the shared `media` table (owner_type +
owner_id) — the most template-consistent choice given the single pgx datasource.

| Method & Path | Auth | Description |
|---|---|---|
| `GET /apis/competition/v1/competitions` | public | List competitions (`?status=draft\|live\|closed`) |
| `GET /apis/competition/v1/competitions/{id}` | public | One competition + its media |
| `GET /apis/competition/v1/categories` | public | List categories |
| `GET /apis/competition/v1/content` | public | Site copy key-value store (hero, footer, …) |
| `POST /apis/competition/v1/admin/competitions` | admin | Create |
| `PUT /apis/competition/v1/admin/competitions/{id}` | admin | Full-field update (all fields incl. slug, category, dates) |
| `DELETE /apis/competition/v1/admin/competitions/{id}` | admin | Delete |
| `POST/PUT/DELETE /apis/competition/v1/admin/categories[/{id}]` | admin | Category CRUD (delete blocks/reassigns while in use) |
| `PUT /apis/competition/v1/admin/content/{key}` | admin | Edit one site-copy value |

Admin routes live under the `/admin/` segment (role-guarded). Money is stored
as integer pence (`ticket_price_pence`). The full-field update validates status
transitions (forward-only draft→live→closed), slug uniqueness, and refuses to
drop `tickets_total` below `tickets_sold`; `tickets_sold` itself is derived and
never writable. **Categories** are first-class (own table, FK from
competitions; deleting one that is still in use returns `409` unless
`?reassign_to=<id>` moves the competitions first). **Site content** is a small
public-read/admin-write key-value store seeded from the previously hardcoded
frontend copy, so the homepage renders its hero/footer/marketing text from the
API.

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
| `PUT /apis/user/v1/admin/users/{id}` | admin | Edit profile (name, email; unique → 409) |
| `POST /apis/user/v1/admin/users/{id}/suspend` \| `/activate` | admin | Toggle `is_active` |

Purchase runs in a transaction: it checks the user is active (suspended → 403),
reads the competition's `ticket_price_pence` from the shared DB, inserts the
ticket rows, and atomically increments the user's `tickets_owned` +
`total_spent_pence`. It does **not** write `competitions.tickets_sold` (owned by
the competition service — a real system would sync that via a JetStream event).
`tickets_owned`/`total_spent_pence` are derived and have no admin setter.
Registration (a public `POST`) is the one deliberate exception to "writes are
admin-only".

### Draw / Winners service (`cmd/draw`)

Owns the `draws` table. Admin creates + runs draws; the public homepage reads
completed results.

| Method & Path | Auth | Description |
|---|---|---|
| `GET /apis/draw/v1/admin/draws` | admin | Searchable, paginated draw list (`?q=&limit=&offset=`) |
| `GET /apis/draw/v1/admin/draws/{id}` | admin | One draw (incl. pending) |
| `POST /apis/draw/v1/admin/draws` | admin | Create a pending draw for a competition |
| `POST /apis/draw/v1/admin/draws/{id}/run` | admin | Run the draw — pick a winner |
| `PUT /apis/draw/v1/admin/draws/{id}` | admin | Edit prize text (void draws frozen) |
| `POST /apis/draw/v1/admin/draws/{id}/void` | admin | Void with a **required** reason |
| `POST /apis/draw/v1/admin/draws/{id}/reassign` | admin | Direct winner reassign (required reason, audited) |
| `GET /apis/draw/v1/draws/{id}` | public | Read a drawn result (pending/void → 404) |
| `GET /apis/draw/v1/winners` | public | Public winners feed (drawn draws + names) |

Running a draw is a single transaction: it guards that the draw is still
`pending` (inside the tx, with a conditional `WHERE ... AND status = 'pending'`
UPDATE so concurrent runs can't both win), reads the competition's tickets from
the shared DB, and picks a uniformly-random winner with **crypto/rand** (not
math/rand). Like the ticket-purchase flow it does not mutate competition-owned
state. The draw service is the sensitive surface, so **every** admin mutation
(create, run, edit, void, reassign) writes an `admin_audit_log` entry with the
acting admin (from JWT claims) and reason. The safe way to change a winner is
void + create + run; direct reassignment is a separate, reason-required,
audited endpoint. The public winners feed replaces the old pattern where the
frontend minted an admin token to read winners — winners are public.

### Gateway service (`cmd/gateway`)

The single public entrypoint. It holds no state — it reverse-proxies (stdlib
`httputil.NewSingleHostReverseProxy`) by the `<servicename>` path segment of
`/apis/<servicename>/...` to that service's upstream base URL (from
`gateway.upstreams.*`, `APP_`-overridable). Unknown service → `404`. Trace
context is propagated to upstreams via an otel-instrumented transport.

**Two-layer (defense-in-depth) role-based admin auth.** A shared HS256
bearer-token middleware (`pkg/middlewares/jwtauth.go`) guards any path matching
`/apis/<svc>/v1/admin/...`:

1. at the **gateway**, before proxying, and
2. **inside each service** (its own admin route group), so a service reached
   directly on its internal port is never unprotected.

The token must be validly signed AND carry a `role` claim of `admin` or
`superadmin`. A missing/invalid/expired token → `401`; a valid token with the
wrong (or absent) role → `403`. The adminauth account-management group
(`/apis/adminauth/v1/admin/accounts`) additionally requires `superadmin`. All
guards read the **same** `jwt.secret` that adminauth signs with. Public reads
and the public POSTs (register, purchase, login/refresh/logout) need no token.
This closes the previous hole where any validly-signed token — including one
forged with the committed default secret from the frontend — was accepted as an
admin.

### Running the full stack locally

One command brings up everything — Postgres, MinIO, migrations, the media
bucket, and all five services behind the gateway:

```sh
cp .env.example .env
docker compose up -d --build
```

- **Only the gateway is published** — `http://localhost:8080`. All public
  traffic goes through `http://localhost:8080/apis/<service>/v1/...`. The four
  domain services (competition, user, draw, media) listen on `:8080` **inside**
  the compose network only and are reachable by service name.
- **Migrations** (all files in `migrations/`, applied with `up`) run in a
  one-shot `migrate/migrate` container after Postgres is healthy and before the
  services start (they wait on `condition: service_completed_successfully`).
  This includes the adminauth schema (000007), categories (000008), user
  `is_active` (000009), site content (000010), draw void (000011) and the audit
  log (000006).
- **The `botb-media` bucket** is created on first run by a one-shot `mc`
  container, so media uploads work immediately.
- The shared `jwt.secret` (`APP_JWT_SECRET`) is defined once in `.env` and
  injected into the gateway, all five services, and adminauth (adminauth signs,
  the rest verify — two-layer role-based admin auth).
- **First superadmin**: set `ADMIN_BOOTSTRAP_EMAIL` / `ADMIN_BOOTSTRAP_PASSWORD`
  in `.env` (defaults `admin@botb.local` / `change-me-now`). adminauth seeds it
  on first boot when no superadmin exists. Log in at `http://localhost:3000/admin`.

Smoke-test through the gateway:

```sh
curl -i http://localhost:8080/apis/competition/v1/competitions                 # public → 200
curl -i -X POST http://localhost:8080/apis/competition/v1/admin/competitions   # no token → 401

# Log in as the bootstrapped superadmin and use the token on an admin route:
TOKEN=$(curl -s -X POST http://localhost:8080/apis/adminauth/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@botb.local","password":"change-me-now"}' | jq -r .access_token)
curl -i http://localhost:8080/apis/user/v1/admin/users -H "Authorization: Bearer $TOKEN"  # → 200
```

Tear down (add `-v` to also drop the Postgres/MinIO volumes):

```sh
docker compose down
```

**Swagger.** Each service serves Swagger UI at `/swagger/` (spec at
`/docs/swagger/swagger.json`) on its internal `:8080`. Because only the gateway
is published, view a service's docs by adding a temporary port mapping (e.g.
`ports: ["8091:8080"]` under `competition:` → open `http://localhost:8091/swagger/`),
or fetch the spec through the compose network:

```sh
docker compose exec gateway wget -qO- http://competition:8080/docs/swagger/swagger.json
```

> **Images build from vendored deps** (`go mod vendor`) so the build needs no
> access to the Go module proxy. Run `go mod vendor` after changing dependencies.
> Each service image is selected by the `SERVICE` build-arg; bases are
> `golang:1.25-alpine` (build) and `alpine:3.20` (runtime).

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
