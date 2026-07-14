# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:0f70d7d828acd8456022127f31975364e58d792999a7e92af6fc972e124bb6b0 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
