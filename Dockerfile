# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:98fc714bfe32e7d3c539d63bda9b9cd089fd699dc3cbd1c534fec3c4deb9ca98 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
