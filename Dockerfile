# syntax=docker/dockerfile:1
FROM golang:1.19.2@sha256:2fddf0539591f8e364c9adb3d495d1ba2ca8a8df420ad23b58e7bcee7986ea6c AS builder
ENV CGO_ENABLED=0
WORKDIR /shortener-bin
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o shortener ./cmd/shortener/

FROM alpine:3.16.2@sha256:1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder ./shortener-bin/shortener .
EXPOSE 8090
CMD ["./shortener"]