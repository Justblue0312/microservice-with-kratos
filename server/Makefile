GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

.PHONY: proto
proto:
	buf generate --template proto/buf.gen.yaml proto

.PHONY: build-hello
build-hello:
	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/hello.exe ./services/hello/cmd/hello

.PHONY: build-goodbye
build-goodbye:
	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/goodbye.exe ./services/goodbye/cmd/goodbye

.PHONY: build-gateway
build-gateway:
	go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/gateway.exe ./gateway/cmd/gateway

.PHONY: build
build: build-hello build-goodbye build-gateway

.PHONY: run-hello
run-hello:
	go run ./services/hello/cmd/hello

.PHONY: run-goodbye
run-goodbye:
	go run ./services/goodbye/cmd/goodbye

.PHONY: run-gateway
run-gateway:
	go run ./gateway/cmd/gateway

.PHONY: generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
all: proto build

.PHONY: clean
clean:
	rm -rf bin/

help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
