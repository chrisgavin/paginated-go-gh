# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:7caba5286b4c3613a337b709c573047d8ae62ee76106647313b61e72b99f20af AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
