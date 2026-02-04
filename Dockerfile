# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:d2e5acc5c29cc331ad5e4f59b09dee0f7d47043b072de0799be63e5c49f95feb AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
