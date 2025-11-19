# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:f60eaa87c79e604967c84d18fd3b151b3ee3f033bcdade4f3494e38411e60963 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
