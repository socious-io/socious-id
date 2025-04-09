# syntax=docker/dockerfile:1.4

########## Base Stage ##########
FROM golang:1.22.5 AS base

WORKDIR /app

# Set for max concurrent downloads
ENV CGO_ENABLED=1
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOMAXPROCS=4
ENV GOMODCACHE=/gomod-cache
ENV GOCACHE=/go-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=$GOMODCACHE,source=~/go/pkg/mod \
    --mount=type=cache,target=$GOCACHE,source=~/.cache/go-build \
    go mod download
COPY . .

#########################
# Test Stage
#########################
FROM base AS test
RUN --mount=type=cache,target=$GOMODCACHE,source=~/go/pkg/mod \
    --mount=type=cache,target=$GOCACHE,source=~/.cache/go-build \
    go mod download
CMD go test -v ./tests -count=1

#########################
# Migration Stage
#########################
FROM base AS migration
CMD go run cmd/migrate/main.go up

#########################
# Build Stage
#########################
FROM base AS builder
RUN --mount=type=cache,target=$GOMODCACHE,source=~/go/pkg/mod \
    --mount=type=cache,target=$GOCACHE,source=~/.cache/go-build \
    go mod download
CMD go build -C cmd/app -trimpath -ldflags="-s -w" -o /app/build .

#########################
# Final Minimal Image
#########################
FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /app/build /app/build
ENTRYPOINT ["/app/build"]