# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:8c5d338aa0da7e8b034efab738460793dc48d2aac8d497c8f8617d4ef964e0a7 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
