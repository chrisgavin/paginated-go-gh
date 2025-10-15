# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:6ea52a02734dd15e943286b048278da1e04eca196a564578d718c7720433dbbe AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
