# syntax=docker/dockerfile:1

FROM golang:1.22.5 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -C cmd/app -o ../../build

EXPOSE 5062
CMD ["./build"]