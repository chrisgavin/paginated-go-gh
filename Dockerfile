# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:a9c4aaca2982a34e9e5b8d5e906f29cf797a893655182bcf8b90d65071a730ea AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
