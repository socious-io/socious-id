# syntax=docker/dockerfile:1.4

########## Base Stage ##########
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Set for max concurrent downloadss
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

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -C cmd/migrate -trimpath -ldflags="-s -w" -o ../../migrate .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -C cmd/app -trimpath -ldflags="-s -w" -o ../../build .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -C cmd/worker -trimpath -ldflags="-s -w" -o ../../worker .

#########################
# Final Minimal Image
#########################
FROM debian:bullseye-slim AS runner
COPY . .
COPY --from=builder /app/build /app/build
COPY --from=builder /app/migrate /app/migrate
ENTRYPOINT ["/app/build"]

FROM debian:bullseye-slim AS worker-runner
COPY . .
COPY --from=builder /app/worker /app/worker
ENTRYPOINT ["/app/worker"]