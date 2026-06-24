# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:b33a0f13c3badc08675233bc9dc064a88124a26e7c1356b1bf9da4a74834ae2b AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
