# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:c0bf2bc2f8e5720aa2e83e45d2085edbf2ad085e2d1a195bb6c3c402350fe661 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
