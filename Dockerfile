# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:8bbd14091f2c61916134fa6aeb8f76b18693fcb29a39ec6d8be9242c0a7e9260 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
