# XXIX Server — Monorepo Architecture Design

> Stack: Go · Chi · PostgreSQL · Redis · Kafka/NATS · Casbin · Protobuf/gRPC (internal)

---

## Repository Layout

```
xxix/
├── gateway/                        # API gateway — single public entry point
├── services/                       # Bounded-context microservices
│   ├── auth/
│   ├── workspace/
│   ├── story/
│   ├── content/
│   └── notification/
├── kit/                            # Shared Go packages (no business logic)
├── infrastructure/                 # Infra-as-code & runtime config
├── proto/                          # Protobuf definitions (inter-service contracts)
├── scripts/                        # Dev/CI helpers
├── .github/
│   └── workflows/
├── go.work                         # Go workspace (links all modules)
├── go.work.sum
├── Makefile
└── README.md
```

---

## Top-level Directories

### `gateway/`

The single public-facing node. Handles TLS termination, routing, rate-limiting, auth-token validation, and fan-out to downstream services.

```
gateway/
├── cmd/
│   └── gateway/
│       └── main.go
├── internal/
│   ├── conf/                       # Config structs (mirrors Kratos pattern)
│   ├── server/                     # HTTP server bootstrap
│   ├── proxy/                      # Route → service forwarding rules
│   └── middleware/                 # Gateway-level middleware
├── configs/
│   ├── config.yaml
│   └── config.prod.yaml
├── Dockerfile
├── go.mod
└── go.sum
```

---

### `kit/`

Pure library modules — zero business logic, zero service imports.

```
kit/
├── apierrors/          # Canonical error codes + HTTP mapping
├── render/             # Content-type-aware response writer
├── serializer/         # JSON/Protobuf codec negotiation
├── middleware/         # Chi middleware: requestid, logger, recover, cors, ratelimit
├── auth/
│   ├── subject/        # Principal extraction from JWT/session
│   ├── claims/         # Token claim types & validation
│   └── ctxkeys/        # Context key constants
├── cache/              # Redis + Ristretto hybrid cache client
├── db/
│   ├── pgx/            # pgxpool factory, health check, query helpers
│   ├── migrate/        # golang-migrate runner
│   └── tx/             # Transaction context helpers
├── messaging/
│   ├── nats/           # NATS JetStream publisher/subscriber wrappers
│   └── kafka/          # Kafka producer/consumer wrappers (sarama)
├── snowid/             # 48-bit distributed ID generator
├── pagination/         # Cursor & offset pagination helpers
├── validator/          # go-playground/validator wrapper
├── telemetry/          # OpenTelemetry tracer + meter bootstrap
├── logger/             # slog wrapper with trace-id injection
└── config/             # Viper-based config loader
```

Each subdirectory is an independent Go package. Services import only what they need.

---

### `proto/`

Single source of truth for all inter-service contracts.

```
proto/
├── auth/v1/
│   └── auth.proto
├── workspace/v1/
│   └── workspace.proto
├── story/v1/
│   └── story.proto
├── content/v1/
│   └── content.proto
├── notification/v1/
│   └── notification.proto
└── buf.yaml
```

Generated Go stubs land in `gen/go/` (gitignored, produced by `make proto`).

---

### `infrastructure/`

```
infrastructure/
├── docker/
│   ├── compose.yaml                # Full local stack
│   ├── compose.dev.yaml            # Override for hot-reload
│   └── compose.test.yaml           # Isolated test DBs
├── terraform/
│   ├── modules/
│   │   ├── gke/
│   │   ├── cloudsql/
│   │   └── redis/
│   └── envs/
│       ├── staging/
│       └── production/
├── kubernetes/
│   ├── base/
│   │   ├── namespace.yaml
│   │   └── configmap.yaml
│   └── overlays/
│       ├── staging/
│       └── production/
├── nginx/
│   └── nginx.conf
├── envoy/
│   └── envoy.yaml
├── prometheus/
│   ├── prometheus.yaml
│   └── rules/
├── grafana/
│   └── dashboards/
└── migrations/                     # Per-service migration files
    ├── auth/
    ├── workspace/
    ├── story/
    ├── content/
    └── notification/
```

---

## Service Template

Every service under `services/<name>/` is an independent Go module that follows this layout:

```
services/<name>/
├── cmd/
│   └── <name>/
│       └── main.go                 # Dependency injection entry point
├── internal/
│   ├── conf/                       # Config structs (env + file)
│   ├── server/                     # HTTP + gRPC server wiring
│   ├── domain/                     # Business entities & interfaces
│   ├── app/                        # Use-cases / application layer
│   ├── repository/                 # DB adapters (implements domain interfaces)
│   ├── handler/                    # HTTP handlers (Chi routes)
│   ├── grpc/                       # gRPC server implementations
│   ├── event/                      # Message broker consumers & publishers
│   └── middleware/                 # Service-local middleware
├── configs/
│   ├── config.yaml
│   └── config.prod.yaml
├── Dockerfile
├── go.mod
└── go.sum
```

### Layer Responsibilities

```
handler / grpc  →  app (use-cases)  →  domain interfaces  ←  repository / event
                         ↑
               kit/* (no business logic)
```

| Layer | Owns | Must NOT import |
|---|---|---|
| `domain/` | Entities, value objects, repository interfaces, domain events | Everything above it |
| `app/` | Use-case orchestration, transaction boundaries | `handler`, `repository` directly |
| `repository/` | SQL queries, pgx, Redis reads/writes | `app`, `handler` |
| `handler/` | HTTP decode → app call → HTTP encode | `repository` directly |
| `event/` | Broker I/O, message schema mapping | `handler` |

---

## Detailed Internal Structure

### `internal/conf/`

```go
// config.go
type Config struct {
    HTTP     HTTPConfig
    GRPC     GRPCConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Broker   BrokerConfig  // NATS or Kafka
    Auth     AuthConfig
}
```

Loaded via `kit/config` (Viper). Config file + env var overrides. No global state — passed through dependency injection.

---

### `internal/server/`

```
server/
├── server.go       # Wires HTTP + gRPC, handles lifecycle
├── http.go         # Chi router setup, middleware chain
└── grpc.go         # gRPC server setup, interceptors
```

```go
// server.go
type Server struct {
    http *http.Server
    grpc *grpc.Server
}

func (s *Server) Start(ctx context.Context) error { ... }
func (s *Server) Stop(ctx context.Context) error  { ... }
```

Graceful shutdown: `Stop` is called on `SIGTERM`/`SIGINT` with a configurable drain timeout.

---

### `internal/domain/`

```
domain/
├── entity.go           # Core structs (no db tags)
├── errors.go           # Domain-scoped sentinel errors
├── repository.go       # Repository interfaces
└── events.go           # Domain event types
```

```go
// repository.go — example (auth service)
type TokenRepository interface {
    Create(ctx context.Context, t *Token) error
    FindByID(ctx context.Context, id snowid.ID) (*Token, error)
    Revoke(ctx context.Context, id snowid.ID) error
    ListActive(ctx context.Context, subjectID snowid.ID) ([]*Token, error)
}
```

---

### `internal/app/`

```
app/
├── <usecase_a>.go
├── <usecase_b>.go
└── service.go          # Aggregates use-cases, holds injected repos
```

Use-cases are plain functions or methods on a `Service` struct. They own transaction boundaries via `kit/db/tx`.

---

### `internal/repository/`

```
repository/
├── postgres/
│   ├── token_repo.go       # Implements domain.TokenRepository
│   └── queries/
│       └── token.sql       # sqlc or raw SQL constants
├── redis/
│   └── token_cache.go
└── repository.go           # Factory / wiring
```

Uses `kit/db/pgx.BaseRepository[T]` (the smart executor pattern) for typed query helpers.

---

### `internal/handler/`

```
handler/
├── router.go       # r.Route("/v1", ...) — mounts all sub-routers
├── <resource>.go   # One file per REST resource
└── dto/
    ├── request.go
    └── response.go
```

```go
// router.go
func New(svc *app.Service, mw ...func(http.Handler) http.Handler) http.Handler {
    r := chi.NewRouter()
    r.Use(mw...)
    r.Route("/v1", func(r chi.Router) {
        r.Mount("/tokens", tokenHandler(svc))
    })
    return r
}
```

Handlers only: decode request → call app → encode response. No business logic.

---

### `internal/event/`

```
event/
├── consumer.go     # Subscribes to broker topics
├── publisher.go    # Publishes domain events
└── handlers/
    └── <topic>.go  # One handler per topic
```

Event handlers call `app` use-cases, same as HTTP handlers. Schema is validated against `proto/` generated types.

---

### `cmd/<name>/main.go` — Dependency Injection

```go
func main() {
    cfg := conf.Load()

    // Infrastructure
    pool    := kit_pgx.NewPool(cfg.Database)
    redis   := kit_cache.NewRedis(cfg.Redis)
    broker  := kit_nats.NewClient(cfg.Broker)
    tracer  := kit_telemetry.NewTracer(cfg.Telemetry)

    // Repositories
    tokenRepo := postgres.NewTokenRepository(pool)

    // Application
    svc := app.NewService(tokenRepo, redis)

    // Transport
    httpHandler := handler.New(svc,
        kit_middleware.RequestID(),
        kit_middleware.Logger(),
        kit_middleware.Recover(),
    )
    srv := server.New(cfg.HTTP, cfg.GRPC, httpHandler)

    // Run
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer stop()

    if err := srv.Start(ctx); err != nil {
        log.Fatal(err)
    }
}
```

Explicit construction, easy to trace. No code generation, no build tags.

---

## Service Inventory

| Service | Responsibilities | Key Topics (NATS/Kafka) |
|---|---|---|
| **auth** | JWT issuance, token lifecycle, Casbin enforcement, session revocation | `auth.token.revoked`, `auth.session.created` |
| **workspace** | Workspace CRUD, member management, role assignment, invite flow | `workspace.member.joined`, `workspace.deleted` |
| **story** | Story projects, outlines, graph canvas, chapter structure | `story.created`, `story.published` |
| **content** | Rich-text documents, versions, diffs, asset uploads | `content.version.saved`, `content.asset.uploaded` |
| **notification** | Fan-out, delivery (email/push/in-app), preference management | `notification.send` (consumer) |

---

## Go Workspace (`go.work`)

```
go 1.23

use (
    ./gateway
    ./services/auth
    ./services/workspace
    ./services/story
    ./services/content
    ./services/notification
    ./kit
)
```

Each module has its own `go.mod`. The workspace links them locally so cross-module changes are reflected immediately without publishing.

---

## `Makefile` — Common Targets

```makefile
.PHONY: proto build test lint migrate-up migrate-down dev

proto:
    buf generate proto/

build:
    go build ./services/... ./gateway/...

test:
    go test -race ./services/... ./gateway/... ./kit/...

lint:
    golangci-lint run ./...

migrate-up:
    go run ./scripts/migrate up --service=$(SVC)

migrate-down:
    go run ./scripts/migrate down --service=$(SVC) --steps=1

dev:
    docker compose -f infrastructure/docker/compose.yaml                    -f infrastructure/docker/compose.dev.yaml up
```

---

## Communication Matrix

```
Client
  │
  ▼
Gateway (Chi, public HTTPS)
  │  JWT validation, routing
  ├──► Auth Service        (HTTP/gRPC)
  ├──► Workspace Service   (HTTP/gRPC)
  ├──► Story Service       (HTTP/gRPC)
  └──► Content Service     (HTTP/gRPC)

Services → Services:
  auth       ──gRPC──►  workspace   (validate membership scopes)
  story      ──gRPC──►  content     (fetch document snapshots)
  any        ──NATS──►  notification (publish send requests)

Async bus (NATS JetStream / Kafka):
  auth       publishes  auth.token.revoked
  workspace  publishes  workspace.member.joined
  story      publishes  story.published
  notification consumes all of the above
```

---

## Naming Conventions

| Concern | Convention |
|---|---|
| Package names | `lowercase`, single word |
| Interfaces | `VerbNoun` (e.g. `TokenRepository`) |
| Constructors | `New<Type>(deps...) *Type` |
| Config structs | `<Name>Config` |
| Context keys | typed private constants in `kit/auth/ctxkeys` |
| SQL files | `queries/<resource>.sql` |
| Proto packages | `xxix.<service>.v1` |
| NATS subjects | `<service>.<resource>.<event>` (e.g. `auth.token.revoked`) |
| Kafka topics | `xxix.<service>.<event>` |

---

## Migration from Kratos Layout

| Kratos | XXIX Monorepo |
|---|---|
| `internal/biz/` | `services/<name>/internal/domain/` + `app/` |
| `internal/data/` | `services/<name>/internal/repository/` |
| `internal/service/` | `services/<name>/internal/handler/` + `grpc/` |
| `internal/server/` | `services/<name>/internal/server/` ✓ kept |
| `internal/conf/` | `services/<name>/internal/conf/` ✓ kept |
| `api/helloworld/v1/` | `proto/<service>/v1/` |
| `third_party/` | dropped — use `go.mod` deps directly |
