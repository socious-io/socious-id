# syntax=docker/dockerfile:1.4

########## Base Stage ##########
FROM golang:1.22.5-alpine AS base
WORKDIR /app

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

########## Runner Stage ##########
FROM base AS runner
WORKDIR /app
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -C cmd/app -trimpath -ldflags="-s -w" -o ../../build .
CMD ["/app/build"]

########## Worker Stage ##########
FROM base AS worker-runner
WORKDIR /app
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -C cmd/worker -trimpath -ldflags="-s -w" -o ../../worker .
CMD ["/app/worker"]