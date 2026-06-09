# XXIX Server — Implementation Guide

> Step-by-step from `kratos new` to a running gateway + hello service.
> Stack: Kratos v2 · Chi · Protobuf · Go Workspace

---

## Before You Start — What NOT to Touch Yet

Two areas intentionally left out of this guide:

**`kit/`** — Don't write shared packages until at least two services actually need the same code. Build it from extraction, not anticipation. Premature kit packages couple services before you understand the real boundaries. Start with code duplicated inside each service; refactor into kit when the pattern stabilises.

**`infrastructure/`** — Docker Compose for local dev is the only thing worth setting up now (Postgres + Redis + NATS). Skip Terraform, Kubernetes overlays, Envoy, and Prometheus until you have something to deploy. Every hour spent on infra config before the first service works is wasted.

Focus order: `proto` → `gateway` → one real service → connect them together → then kit/infra.

---

## Prerequisites

```bash
# Go 1.23+
go version

# Kratos CLI
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
kratos --version

# Buf (protobuf toolchain)
go install github.com/bufbuild/buf/cmd/buf@latest

# protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
```

---

## Step 1 — Scaffold with Kratos, Then Reshape

### 1.1 Bootstrap the root

```bash
mkdir xxix && cd xxix
kratos new xxix          # generates a single-service template
```

Kratos creates this:
```
xxix/
├── api/helloworld/v1/
├── cmd/xxix/
├── configs/
├── internal/{biz,data,service,server,conf}/
└── third_party/
```

### 1.2 Tear down what you don't need

```bash
# Remove Kratos-specific internal layers we're replacing
rm -rf internal/biz internal/data internal/service
rm -rf api/                  # proto definitions move to proto/
rm -rf third_party/          # dependencies handled by go.mod
```

### 1.3 Reshape into the monorepo layout

```bash
# Promote the remaining internals into the gateway module
mkdir -p gateway/cmd/gateway
mkdir -p gateway/internal/{conf,server,proxy,middleware}
mkdir -p gateway/configs

mv cmd/xxix/main.go           gateway/cmd/gateway/main.go
mv internal/server/           gateway/internal/server/
mv internal/conf/             gateway/internal/conf/
mv configs/                   gateway/configs/

# Root-level structure
mkdir -p services/{auth,workspace,story,content,notification}
mkdir -p proto/{auth,workspace,story,content,notification,helloworld}
mkdir -p kit scripts infrastructure/docker
```

### 1.4 Root go.work

```bash
# In xxix/ root
go work init
go work use ./gateway
```

You'll add each service module here as you create them.

### 1.5 Gateway go.mod

```bash
cd gateway
go mod init github.com/yourorg/xxix/gateway

go get github.com/go-kratos/kratos/v2@latest
go get github.com/go-chi/chi/v5@latest
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf@latest
```

---

## Step 2 — Proto Setup

### 2.1 buf.yaml at repo root

```yaml
# proto/buf.yaml
version: v2
modules:
  - path: .
deps:
  - buf.build/googleapis/googleapis
  - buf.build/grpc-ecosystem/grpc-gateway
lint:
  use:
    - DEFAULT
breaking:
  use:
    - FILE
```

### 2.2 buf.gen.yaml at repo root

```yaml
# buf.gen.yaml  (sits next to proto/)
version: v2
plugins:
  - plugin: go
    out: gen/go
    opt: paths=source_relative

  - plugin: go-grpc
    out: gen/go
    opt:
      - paths=source_relative
      - require_unimplemented_servers=false

  - plugin: go-http          # Kratos HTTP binding generator
    out: gen/go
    opt: paths=source_relative

  - plugin: go-errors        # Kratos error enum generator
    out: gen/go
    opt: paths=source_relative
```

### 2.3 Helloworld proto (sample service)

```protobuf
// proto/helloworld/v1/helloworld.proto
syntax = "proto3";

package helloworld.v1;

option go_package = "github.com/yourorg/xxix/gen/go/helloworld/v1;helloworldv1";

import "google/api/annotations.proto";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {
    option (google.api.http) = {
      get: "/v1/hello"
    };
  }
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

### 2.4 Generate

```bash
# From repo root
buf generate proto

# Adds generated files under gen/go/
# Add gen/go/ to .gitignore OR commit it — pick one policy and stick to it.
# Recommendation: commit gen/ so CI doesn't need buf installed.
```

---

## Step 3 — Hello Service (Sample Microservice)

This is the template every future service copies. Build it correctly once.

### 3.1 Module init

```bash
cd services/hello        # new folder, not a Kratos template
go mod init github.com/yourorg/xxix/services/hello

go get github.com/go-kratos/kratos/v2@latest
go get github.com/go-chi/chi/v5@latest
go get google.golang.org/grpc@latest
```

```bash
# Back at workspace root, register it
cd ..
go work use ./services/hello
```

### 3.2 Service layout

```
services/hello/
├── cmd/hello/
│   └── main.go
├── internal/
│   ├── conf/
│   │   └── config.go
│   ├── server/
│   │   ├── server.go
│   │   ├── http.go
│   │   └── grpc.go
│   ├── domain/
│   │   └── greeter.go          # interface + entity
│   ├── app/
│   │   └── greeter.go          # use-case
│   ├── handler/
│   │   └── greeter.go          # HTTP handler (Chi)
│   └── grpc/
│       └── greeter.go          # gRPC server (implements proto)
├── configs/
│   └── config.yaml
├── go.mod
└── go.sum
```

### 3.3 Config

```go
// internal/conf/config.go
package conf

type Config struct {
    HTTP HTTPConfig
    GRPC GRPCConfig
}

type HTTPConfig struct {
    Addr string // ":8081"
}

type GRPCConfig struct {
    Addr string // ":9081"
}
```

```yaml
# configs/config.yaml
http:
  addr: ":8081"
grpc:
  addr: ":9081"
```

### 3.4 Domain layer

```go
// internal/domain/greeter.go
package domain

import "context"

type GreetRequest struct {
    Name string
}

type GreetReply struct {
    Message string
}

// GreeterService is the interface the app layer implements.
// Repository interfaces go here too when you have DB access.
type Greeter interface {
    Greet(ctx context.Context, req *GreetRequest) (*GreetReply, error)
}
```

### 3.5 Application layer

```go
// internal/app/greeter.go
package app

import (
    "context"
    "fmt"

    "github.com/yourorg/xxix/services/hello/internal/domain"
)

type GreeterService struct{}

func NewGreeterService() *GreeterService {
    return &GreeterService{}
}

func (s *GreeterService) Greet(_ context.Context, req *domain.GreetRequest) (*domain.GreetReply, error) {
    return &domain.GreetReply{Message: fmt.Sprintf("Hello, %s!", req.Name)}, nil
}
```

### 3.6 HTTP handler (Chi)

```go
// internal/handler/greeter.go
package handler

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/yourorg/xxix/services/hello/internal/domain"
)

type GreeterHandler struct {
    svc domain.Greeter
}

func NewGreeterHandler(svc domain.Greeter) *GreeterHandler {
    return &GreeterHandler{svc: svc}
}

func (h *GreeterHandler) Routes() func(r chi.Router) {
    return func(r chi.Router) {
        r.Get("/hello", h.sayHello)
    }
}

func (h *GreeterHandler) sayHello(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "world"
    }
    reply, err := h.svc.Greet(r.Context(), &domain.GreetRequest{Name: name})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reply)
}
```

### 3.7 gRPC server (implements generated proto)

```go
// internal/grpc/greeter.go
package grpcserver

import (
    "context"

    helloworldv1 "github.com/yourorg/xxix/gen/go/helloworld/v1"
    "github.com/yourorg/xxix/services/hello/internal/domain"
    "google.golang.org/grpc"
)

type GreeterServer struct {
    helloworldv1.UnimplementedGreeterServer
    svc domain.Greeter
}

func NewGreeterServer(svc domain.Greeter) *GreeterServer {
    return &GreeterServer{svc: svc}
}

func (s *GreeterServer) SayHello(ctx context.Context, req *helloworldv1.HelloRequest) (*helloworldv1.HelloReply, error) {
    reply, err := s.svc.Greet(ctx, &domain.GreetRequest{Name: req.Name})
    if err != nil {
        return nil, err
    }
    return &helloworldv1.HelloReply{Message: reply.Message}, nil
}

func (s *GreeterServer) Register(srv *grpc.Server) {
    helloworldv1.RegisterGreeterServer(srv, s)
}
```

### 3.8 Server wiring (Kratos transport)

```go
// internal/server/http.go
package server

import (
    "net/http"

    kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
    "github.com/go-chi/chi/v5"
    chiMiddleware "github.com/go-chi/chi/v5/middleware"
    "github.com/yourorg/xxix/services/hello/internal/conf"
    "github.com/yourorg/xxix/services/hello/internal/handler"
)

func NewHTTPServer(cfg *conf.Config, greeter *handler.GreeterHandler) *kratoshttp.Server {
    // Build a Chi router
    r := chi.NewRouter()
    r.Use(chiMiddleware.RequestID)
    r.Use(chiMiddleware.Logger)
    r.Use(chiMiddleware.Recoverer)

    r.Route("/v1", greeter.Routes())

    // Wrap Chi inside a Kratos HTTP server for lifecycle management
    srv := kratoshttp.NewServer(
        kratoshttp.Address(cfg.HTTP.Addr),
    )
    srv.HandlePrefix("/", r)
    return srv
}
```

```go
// internal/server/grpc.go
package server

import (
    kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/yourorg/xxix/services/hello/internal/conf"
    grpcserver "github.com/yourorg/xxix/services/hello/internal/grpc"
)

func NewGRPCServer(cfg *conf.Config, greeter *grpcserver.GreeterServer) *kratosgrpc.Server {
    srv := kratosgrpc.NewServer(
        kratosgrpc.Address(cfg.GRPC.Addr),
    )
    greeter.Register(srv.Server) // srv.Server is the underlying *grpc.Server
    return srv
}
```

```go
// internal/server/server.go
package server

import (
    "github.com/go-kratos/kratos/v2"
    kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
    kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func NewApp(http *kratoshttp.Server, grpc *kratosgrpc.Server) *kratos.App {
    return kratos.New(
        kratos.Name("hello"),
        kratos.Server(http, grpc),
    )
}
```

### 3.9 Dependency Injection (in main.go)

```go
// cmd/hello/main.go
package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/justblue/luoye/services/hello/internal/conf"
	grpcserver "github.com/justblue/luoye/services/hello/internal/grpchandler"
	"github.com/justblue/luoye/services/hello/internal/httphandler"
	"github.com/justblue/luoye/services/hello/internal/server"
	"github.com/justblue/luoye/services/hello/internal/usecase"
)

func initApp(cfg *conf.Config) (*kratos.App, error) {
	svc := usecase.NewGreeterService()
	httpServer := server.NewHTTPServer(cfg, httphandler.NewGreeterHandler(svc))
	grpcServer := server.NewGRPCServer(cfg, grpcserver.NewGreeterServer(svc))
	return server.NewApp(httpServer, grpcServer), nil
}

func main() {
	cfg := &conf.Config{
		HTTP: conf.HTTPConfig{Addr: ":8081"},
		GRPC: conf.GRPCConfig{Addr: ":9081"},
	}
	app, err := initApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

```bash
cd services/hello/cmd/hello
go run .    # server starts
```

---
## Step 4 — Gateway

The gateway is a Kratos app that **does not own business logic** — it validates tokens, maps routes to upstream services, and handles cross-cutting concerns.

### 4.1 Key mental model (Kratos gateway integration)

Kratos gives you transport servers (`kratoshttp.Server`, `kratosgrpc.Server`) that implement `transport.Server`. The gateway uses **one** Kratos app that boots both an inbound HTTP server (public) and outbound gRPC clients (to each service).

```
Internet → Gateway HTTP (Chi) → gRPC client → Hello Service gRPC
```

The gateway never imports a service's `internal/` packages. It only imports the generated proto stubs from `gen/go/`.

**Three things Kratos does for you in the gateway that you must not reinvent:**
1. `kratos.App` — `Run()` + `Stop()` with graceful shutdown already handled.
2. `kratoshttp.Server` + `kratosgrpc.Server` — both implement `transport.Server` interface; pass them to `kratos.Server(...)`.
3. Middleware chain — Kratos middleware runs at the transport layer (before your handler); use it for tracing, logging, auth.

### 4.2 Gateway module

```bash
cd gateway
go mod init github.com/yourorg/xxix/gateway

go get github.com/go-kratos/kratos/v2@latest
go get github.com/go-chi/chi/v5@latest
go get google.golang.org/grpc@latest

```bash
cd ..
go work use ./gateway
```

### 4.3 Gateway layout (filled in)

```
gateway/
├── cmd/gateway/
│   └── main.go
├── internal/
│   ├── conf/
│   │   └── config.go           # upstream addresses live here
│   ├── server/
│   │   ├── server.go           # kratos.App factory
│   │   └── http.go             # Chi router + proxy routes
│   ├── proxy/
│   │   └── hello.go            # gRPC client → hello service
│   └── middleware/
│       └── auth.go             # JWT validation (stub for now)
├── configs/
│   └── config.yaml
├── go.mod
└── go.sum
```

### 4.4 Config

```go
// internal/conf/config.go
package conf

type Config struct {
    HTTP      HTTPConfig
    Upstreams UpstreamsConfig
}

type HTTPConfig struct {
    Addr string // ":8080"
}

type UpstreamsConfig struct {
    Hello string // "localhost:9081"
}
```

```yaml
# configs/config.yaml
http:
  addr: ":8080"
upstreams:
  hello: "localhost:9081"
```

### 4.5 gRPC proxy client

```go
// internal/proxy/hello.go
package proxy

import (
    helloworldv1 "github.com/yourorg/xxix/gen/go/helloworld/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type HelloProxy struct {
    client helloworldv1.GreeterClient
}

func NewHelloProxy(addr string) (*HelloProxy, error) {
    conn, err := grpc.NewClient(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, err
    }
    return &HelloProxy{client: helloworldv1.NewGreeterClient(conn)}, nil
}

func (p *HelloProxy) Client() helloworldv1.GreeterClient {
    return p.client
}
```

### 4.6 HTTP server (gateway Chi router)

```go
// internal/server/http.go
package server

import (
    "encoding/json"
    "net/http"

    kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
    "github.com/go-chi/chi/v5"
    chiMiddleware "github.com/go-chi/chi/v5/middleware"
    helloworldv1 "github.com/yourorg/xxix/gen/go/helloworld/v1"
    "github.com/yourorg/xxix/gateway/internal/conf"
    "github.com/yourorg/xxix/gateway/internal/proxy"
)

func NewHTTPServer(cfg *conf.Config, hello *proxy.HelloProxy) *kratoshttp.Server {
    r := chi.NewRouter()
    r.Use(chiMiddleware.RequestID)
    r.Use(chiMiddleware.Logger)
    r.Use(chiMiddleware.Recoverer)

    // Route: GET /v1/hello → hello service
    r.Get("/v1/hello", func(w http.ResponseWriter, r *http.Request) {
        name := r.URL.Query().Get("name")
        reply, err := hello.Client().SayHello(r.Context(), &helloworldv1.HelloRequest{Name: name})
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadGateway)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(reply)
    })

    srv := kratoshttp.NewServer(kratoshttp.Address(cfg.HTTP.Addr))
    srv.HandlePrefix("/", r)
    return srv
}
```

### 4.7 App factory

```go
// internal/server/server.go
package server

import (
    "github.com/go-kratos/kratos/v2"
    kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func NewApp(http *kratoshttp.Server) *kratos.App {
    return kratos.New(
        kratos.Name("gateway"),
        kratos.Server(http),
    )
}
```

## Step 5 — Smoke Test

```bash
# Terminal 1 — start hello service
cd services/hello/cmd/hello && go run .

# Terminal 2 — start gateway
cd gateway/cmd/gateway && go run .

# Terminal 3 — test
curl "http://localhost:8080/v1/hello?name=XXIX"
# {"message":"Hello, XXIX!"}

# Direct hello service (bypassing gateway)
curl "http://localhost:8081/v1/hello?name=direct"
# {"message":"Hello, direct!"}
```

---

## Step 6 — Makefile Targets (Root)

```makefile
.PHONY: proto run-hello run-gateway

proto:
	buf generate proto

run-hello:
	go run ./services/hello/cmd/hello/...

run-gateway:
	go run ./gateway/cmd/gateway/...
```

---

## Key Points: Kratos + Services Integration

These are the things that will trip you up if you miss them.

**1. `kratos.App` owns the process lifecycle — don't fight it.**
Call `app.Run()` and let it block. Graceful shutdown on `SIGTERM`/`SIGINT` is already handled. If you try to manage `os.Signal` yourself alongside Kratos, you'll get double-handling.

**2. Transport servers implement `transport.Server` — that's the only contract.**
`kratoshttp.Server` and `kratosgrpc.Server` both implement this interface. Pass them to `kratos.Server(...)`. You can mix-and-match or add a custom server (e.g. a NATS consumer that implements `Start(ctx)/Stop(ctx)`) and Kratos will lifecycle-manage it too.

**3. Chi inside Kratos HTTP: use `srv.HandlePrefix("/", chiRouter)`.**
Kratos HTTP server is not a raw `net/http` server — it wraps one. `HandlePrefix` mounts your Chi router at a path prefix without losing Chi's request context. Don't call `http.Handle` or `http.ListenAndServe` separately.

**4. The gateway imports proto stubs, never service internals.**
`gen/go/helloworld/v1` is the contract. The gateway never touches `services/hello/internal/`. This is what makes services independently deployable. If you find yourself wanting to import a service's internal package from the gateway, that's a signal a proto message is missing.

**5. gRPC clients in the gateway use `grpc.NewClient`, not `grpc.Dial`.**
`grpc.Dial` is deprecated as of Go gRPC v1.60. Use `grpc.NewClient` — it doesn't establish the connection immediately (lazy connect), which is correct behaviour for a gateway that starts before all upstreams are ready.

**9. Kratos middleware runs at transport layer, not Chi middleware layer.**
Kratos has its own middleware chain (`kratoshttp.Middleware(...)`). Use it for cross-cutting concerns that need Kratos context (tracing, operation name extraction). Use Chi middleware for HTTP-specific concerns (request ID, CORS, body limit). Don't duplicate them.

**10. `go.work` means no `replace` directives needed.**
When services import `gen/go/` or shared packages, Go workspace resolves them locally. The moment you publish and stop using `go.work` (e.g. in CI), you'll need real module versions. Keep generated code in `gen/` at the root and make sure `go.work use ./` includes it, or put `gen/` inside its own `go.mod`.

---

## Current State After These Steps

```
server/                     # Go workspace root
├── gateway/                ✓ running on :8080
├── services/hello/         ✓ running on :8081 (HTTP) + :9081 (gRPC)
├── proto/helloworld/       ✓ source of truth
├── gen/go/helloworld/      ✓ generated stubs
├── go.work                 ✓ links gateway + hello
├── kit/                    — empty, fill by extraction later
├── scripts/                - scaffold.py etc.
└── Makefile
```

> Note: `server/` is the Go workspace root. All `go` and `buf` commands in this guide are run from `server/`.

Next service (`auth`) follows the identical template: copy `services/hello/`, rename packages, replace the domain layer, add DB/Redis dependencies, register in `go.work`, add a gRPC proxy in `gateway/internal/proxy/auth.go`.
