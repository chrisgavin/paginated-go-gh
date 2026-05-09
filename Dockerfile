# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:2981696eed011d747340d7252620932677929cce7d2d539602f56a8d7e9b660b AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
