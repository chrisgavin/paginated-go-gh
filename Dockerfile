# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:516827db2015144cf91e042d1b6a3aca574d013a4705a6fdc4330444d47169d5 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
