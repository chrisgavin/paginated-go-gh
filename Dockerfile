# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:b39810f6440772ab1ddaf193aa0c2a2bbddebf7a877f127c113b103e48fd8139 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
