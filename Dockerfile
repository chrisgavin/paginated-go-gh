# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:13605dbaf3aff39741644cc3a6ec74cac494f955d1401ffee49b55032fa8a626 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
