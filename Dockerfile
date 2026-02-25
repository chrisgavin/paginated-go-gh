# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:b75179794e029c128d4496f695325a4c23b29986574ad13dd52a0d3ee9f72a6f AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
