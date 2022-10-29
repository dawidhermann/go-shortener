# syntax=docker/dockerfile:1
FROM golang:1.18.3@sha256:5417b4917fa7ed3ad2678a3ce6378a00c95bfd430c2ffa39936fce55130b5f2c AS builder
ENV CGO_ENABLED=0
WORKDIR /shortener-bin
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o shortener ./cmd/shortener/

FROM alpine:3.16.0@sha256:4ff3ca91275773af45cb4b0834e12b7eb47d1c18f770a0b151381cd227f4c253
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder ./shortener-bin/shortener .
EXPOSE 8090
CMD ["./shortener"]