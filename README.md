# XXIX Server — Monorepo Microservices

> Stack: Go 1.26 · Kratos v2 · Chi · Protobuf/gRPC · Go Workspace

## Architecture

```
Client
  │
  ▼
Gateway (Chi, :8080)
  │
  ├──► Hello Service   (gRPC :9081, HTTP :8081)
  │
  └──► Goodbye Service (gRPC :9082, HTTP :8082)

Async bus (NATS JetStream):
  goodbye ──publishes──► goodbye.said ──consumed by──► hello
```

- **Gateway** — single public entry point. Routes HTTP requests to downstream services via gRPC. No business logic.
- **Services** — independent Go modules, each with its own `go.mod`, HTTP + gRPC servers, domain logic.
- **Proto** — single source of truth for inter-service contracts. Generated stubs in `gen/go/`.

## Project Layout

```
server/
├── bin/                        # Compiled binaries
├── gateway/                    # API gateway module
│   ├── cmd/gateway/main.go
│   └── internal/
│       ├── conf/               # Config structs
│       ├── server/             # Kratos HTTP server + Chi router
│       └── proxy/              # gRPC clients → upstream services
├── cmd/provision/              # NATS stream provisioner tool
├── services/
│   ├── hello/                  # Greeter service — consumes goodbye.said via NATS
│   │   ├── cmd/hello/
│   │   │   ├── main.go
│   │   │   ├── wire.go
│   │   │   └── wire_gen.go
│   │   └── internal/
│   │       ├── conf/
│   │       ├── server/         # Kratos HTTP + gRPC server wiring
│   │       ├── domain/         # Entities & interfaces
│   │       ├── usecase/        # Business logic
│   │       ├── grpchandler/    # gRPC handlers (also serves HTTP via Kratos)
│   │       ├── grpcclient/     # Outbound gRPC client (→ goodbye)
│   │       └── event/          # NATS consumer
│   └── goodbye/                # Goodbye service — publishes events, internal async worker
│       ├── cmd/goodbye/
│       │   ├── main.go
│       │   ├── wire.go
│       │   └── wire_gen.go
│       └── internal/
│           ├── conf/
│           ├── server/         # Kratos HTTP + gRPC + Asynq server wiring
│           ├── domain/         # Entities & interfaces
│           ├── usecase/        # Business logic
│           ├── grpchandler/    # gRPC handlers
│           ├── event/          # NATS publisher + Asynq publisher
│           └── worker/         # Async task processor (internal Asynq handler)
├── proto/                      # Protobuf definitions
├── gen/go/                     # Generated Go stubs
├── scripts/                    # Dev helpers
├── kit/                        # Shared packages
│   └── messaging/
│       ├── nats/               # NATS JetStream client wrapper
│       └── asynq/              # Asynq task queue client/server + task types
├── infrastructure/             # Docker Compose, K8s, Terraform
│   └── docker/compose.yaml     # Local dev stack (NATS, Redis, services)
├── go.work                     # Go workspace
└── Makefile
```

## Quick Start

```bash
cd server

# Install tools (one-time)
make init

# Generate proto stubs
make proto

# Build all binaries
make build

# Start infrastructure (NATS + Redis)
docker compose -f infrastructure/docker/compose.yaml up -d nats redis

# Provision NATS streams (one-time per deploy)
go run ./cmd/provision -url nats://localhost:4222

# Terminal 1 — goodbye service (gRPC :9082, HTTP :8082)
make run-goodbye

# Terminal 2 — hello service (gRPC :9081, HTTP :8081)
make run-hello

# Terminal 3 — gateway (:8080)
make run-gateway

# Test
curl "http://localhost:8080/v1/hello?name=XXIX"
# → {"message":"Hello, XXIX! Goodbye, XXIX!"}
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make proto` | Generate protobuf stubs via Buf |
| `make build` | Build all binaries → `bin/` |
| `make build-hello` | Build hello service only |
| `make build-goodbye` | Build goodbye service only |
| `make build-gateway` | Build gateway only |
| `make run-hello` | Run hello service |
| `make run-goodbye` | Run goodbye service |
| `make run-gateway` | Run gateway |
| `make clean` | Remove `bin/` |
| `make generate` | `go generate ./...` + `go mod tidy` |
| `make all` | `proto` + `build` |

## Service Template

Every service under `services/<name>/` follows this structure:

```
services/<name>/
├── cmd/<name>/
│   ├── main.go               # Entry point — loads config, calls initApp
│   ├── wire.go               # Wire DI declarations
│   └── wire_gen.go           # Generated Wire code (go generate)
├── internal/
│   ├── conf/config.go        # Config structs
│   ├── server/
│   │   ├── server.go         # Kratos App factory
│   │   ├── http.go           # Kratos HTTP server (Chi or Kratos routes)
│   │   └── grpc.go           # gRPC server
│   ├── domain/               # Entities & repository interfaces
│   ├── usecase/              # Business logic
│   ├── grpchandler/          # gRPC handlers (also serves HTTP via Kratos transcoding)
│   ├── event/                # NATS publishers/consumers
│   ├── grpcclient/           # Outbound gRPC clients (if calling other services)
│   └── worker/               # Async task processor (optional, goodbye only)
├── config/config.yaml
├── go.mod
└── go.sum
```

Each layer follows a strict dependency direction:

```
grpchandler / event → usecase → domain interfaces ← grpcclient / repository
                              ↑
                    kit/* (no business logic)
```

### Adding a New Service

```bash
# Copy an existing service as template
cp -r services/hello services/<service_name>

# Rename packages and register in workspace
go work use ./services/<service_name>
cd services/<service_name> && go mod tidy

# Define proto in proto/<service_name>/v1/<service_name>.proto
# Add gRPC proxy in gateway/internal/proxy/<service_name>.go

# After changing wire.go
go generate ./services/<service_name>/...
```

## Key Design Decisions

| Decision | Choice | Why |
|----------|--------|-----|
| DI approach | **Google Wire** | Compile-time DI, no runtime reflection. Clean separation between wiring and logic. |
| HTTP router | **Chi** inside Kratos | Chi middleware ecosystem, clean route grouping, Kratos handles lifecycle. |
| Gateway strategy | **gRPC client → service** | Type-safe, performant, single proto contract. |
| Proto generation | **Buf** (local plugins) | Modern, fast, dependency management built-in. |
| Module system | **Go workspace** | Cross-module changes without `replace` directives. |

## Communication Matrix

```
Gateway (Chi, :8080)
  │
  ├──► Hello Service   (gRPC :9081)   — greeter endpoint
  └──► Goodbye Service (gRPC :9082)   — goodbye endpoint

Services → Services:
  hello  ──gRPC──►  goodbye   (SayGoodbye called within Greet)

Async bus (NATS JetStream):
  goodbye  ──publishes──►  goodbye.said  ──consumed by──►  hello

Async task queue (Asynq/Redis, internal to goodbye):
  goodbye  ──enqueues──►  goodbye:said   ──processed by──►  goodbye worker
```
