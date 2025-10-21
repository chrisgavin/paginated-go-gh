# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:0d8c14c93010e9eab98599f1ddd37e80b8fd39e9c662d670c4e4d9d0b101831d AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
