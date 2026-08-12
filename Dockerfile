# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:5822931cf78fe98a97edcf73a0c54c29fa2386b99c8136468e274ae9fab8cfba AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
