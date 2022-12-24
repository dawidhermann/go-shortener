# syntax=docker/dockerfile:1
FROM golang:1.19.3@sha256:24e286ca5b48c690f29266a7086e1b7f77a4ddc1a47f6f8bf55d4b736eee073e AS builder
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