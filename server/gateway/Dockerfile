FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o /gateway ./gateway/cmd/gateway

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /gateway .
COPY gateway/configs/config.yaml configs/config.yaml
EXPOSE 8080
CMD ["./gateway"]
