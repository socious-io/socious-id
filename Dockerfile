# syntax=docker/dockerfile:1.4

########## Base Stage ##########
FROM golang:1.22.5 AS base

WORKDIR /app

# Set for max concurrent downloads
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOMAXPROCS=4
ENV GOMODCACHE=/go/pkg/mod
ENV GOCACHE=/root/.cache/go-build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .

#########################
# Test Stage
#########################
FROM base AS test
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go test -v ./tests

#########################
# Migration Stage
#########################
FROM base AS migration
CMD go run cmd/migrate/main.go up

#########################
# Build Stage
#########################
FROM base AS builder
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -C cmd/app -trimpath -ldflags="-s -w" -o ../../build .

#########################
# Final Minimal Image
#########################
FROM debian:bullseye-slim
COPY . .
COPY --from=builder /app/build /app/build
ENTRYPOINT ["/app/build"]