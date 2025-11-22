# syntax=docker/dockerfile:1
# check=error=true;experimental=all

FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY ./go.mod  ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /api-gateway ./cmd

FROM busybox:1.36.1 AS release

WORKDIR /
COPY --from=builder /app/config/config.yaml /config/config.yaml
COPY --from=builder /api-gateway /api-gateway

ENTRYPOINT ["/api-gateway"]