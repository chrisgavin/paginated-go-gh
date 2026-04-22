# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:1e598ea5752ae26c093b746fd73c5095af97d6f2d679c43e83e0eac484a33dc3 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
