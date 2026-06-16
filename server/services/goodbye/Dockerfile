FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o /goodbye ./services/goodbye/cmd/goodbye

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /goodbye .
COPY services/goodbye/config/config.yaml config/config.yaml
EXPOSE 8082 9082
CMD ["./goodbye"]
