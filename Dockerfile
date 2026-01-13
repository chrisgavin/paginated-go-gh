# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:581c059c1d53b96a6db51a2f7a0eb943491ba1da50f1e43700db0cc325618096 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
