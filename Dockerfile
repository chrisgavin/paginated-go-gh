# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:f7159064a17ccc65d0e10e342ae8783026182704bf4af8f6df8d5ba9af2be2fd AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
