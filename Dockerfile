# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:11fd8f7f63db3b6fb198797042ba4c40a4a34dc83325d3328ca3bc4bb7726786 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
