# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:8f4cb3b8d3fd8c3e6eccfde0fcf54e8cea74fbb04cea961a92ee1a913d22cb17 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
