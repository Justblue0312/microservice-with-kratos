FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o /app/worker ./services/worker/cmd/worker

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/worker .
COPY services/worker/config/config.yaml config/config.yaml
EXPOSE 8083
CMD ["./worker"]
