# syntax=docker/dockerfile:1.4

########## Base Stage ##########
FROM golang:1.22.5 AS base

WORKDIR /app

# Set for max concurrent downloads
ENV CGO_ENABLED=1
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOMAXPROCS=4
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    go mod download
COPY . .

#########################
# Test Stage
#########################
FROM base AS test
WORKDIR /app
COPY . .
COPY --from=base /gomod-cache /gomod-cache
COPY --from=base /go-cache /go-cache
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
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -C cmd/app -trimpath -ldflags="-s -w" -o /app/build .

#########################
# Final Minimal Image
#########################
FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /app/build /app/build
ENTRYPOINT ["/app/build"]