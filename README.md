# XXIX Server — Monorepo Microservices

> Stack: Go 1.26 · Kratos v2 · Chi · Protobuf/gRPC · Go Workspace

## Architecture

```
Client
  │
  ▼
Gateway (Chi, port :8080)
  │  ┌──────────────────────────────┐
  ├──► Hello Service (gRPC :9081)   │  HTTP :8081
  ├──► Auth Service   (gRPC :9082)  │  HTTP :8082
  ├──► Story Service  (gRPC :9083)  │  HTTP :8083
  └──► …                             │
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
├── services/
│   └── hello/                  # Sample service (template for others)
│       ├── cmd/hello/main.go
│       └── internal/
│           ├── conf/
│           ├── server/         # HTTP (Chi) + gRPC server wiring
│           ├── domain/         # Entities & interfaces
│           ├── usecase/        # Application logic
│           ├── httphandler/    # Chi HTTP handlers
│           └── grpchandler/    # gRPC server implementations
├── proto/                      # Protobuf definitions
├── gen/go/                     # Generated Go stubs
├── scripts/                    # Dev helpers
├── kit/                        # Shared packages (extract when needed)
├── infrastructure/             # Docker Compose, K8s, Terraform
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

# Start hello service (Terminal 1)
make run-hello

# Start gateway (Terminal 2)
make run-gateway

# Test (Terminal 3)
curl "http://localhost:8080/v1/hello?name=XXIX"
# → {"message":"Hello, XXIX!"}

curl "http://localhost:8081/v1/hello?name=direct"
# → {"Message":"Hello, direct!"}
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make proto` | Generate protobuf stubs via Buf |
| `make build` | Build all binaries → `bin/` |
| `make build-hello` | Build hello service only |
| `make build-gateway` | Build gateway only |
| `make run-hello` | Run hello service |
| `make run-gateway` | Run gateway |
| `make clean` | Remove `bin/` |
| `make generate` | `go generate ./...` + `go mod tidy` |
| `make all` | `proto` + `build` |

## Service Template

Every service under `services/<name>/` follows the same structure:

```
services/<name>/
├── cmd/<name>/
│   └── main.go               # Dependency injection entry point
├── internal/
│   ├── conf/config.go        # Config structs
│   ├── server/
│   │   ├── server.go         # Kratos App factory
│   │   ├── http.go           # Chi HTTP server
│   │   └── grpc.go           # gRPC server
│   ├── domain/               # Entities, interfaces
│   ├── usecase/              # Business logic
│   ├── httphandler/          # Chi HTTP handlers
│   └── grpchandler/          # gRPC handlers (proto impl)
├── config/config.yaml
├── go.mod
└── go.sum
```

### Adding a New Service

```bash
python scripts/scaffold.py <service_name>

# Then:
go work use ./services/<service_name>
cd services/<service_name> && go mod tidy
# Add gRPC proxy in gateway/internal/proxy/<service_name>.go
# Define proto in proto/<service_name>/<service_name>.proto
```

## Key Design Decisions

| Decision | Choice | Why |
|----------|--------|-----|
| DI approach | **Manual** (no wire) | Shallow dep graphs (depth ≤ 3). Simpler, no extra build step, easier to debug. |
| HTTP router | **Chi** inside Kratos | Chi middleware ecosystem, clean route grouping, Kratos handles lifecycle. |
| Gateway strategy | **gRPC client → service** | Type-safe, performant, single proto contract. |
| Proto generation | **Buf** (local plugins) | Modern, fast, dependency management built-in. |
| Module system | **Go workspace** | Cross-module changes without `replace` directives. |

## Communication Matrix

```
Gateway (Chi, :8080)
  │  JWT validation, routing, rate-limiting
  ├──► Auth Service        (gRPC)
  ├──► Workspace Service   (gRPC)
  ├──► Story Service       (gRPC)
  └──► Content Service     (gRPC)

Services → Services via gRPC
Async bus: NATS JetStream / Kafka
```
