FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o /hello ./services/hello/cmd/hello

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /hello .
COPY services/hello/config/config.yaml config/config.yaml
EXPOSE 8081 9081
CMD ["./hello"]
