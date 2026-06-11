# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:d47ca13cd596f3a338c1be5f79af628f42bedcf89455266211a9ab4f95da2828 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
