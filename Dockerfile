# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:1a13867f40c6ba8c3e7c07ab5fa921cb58cf8402bcbe452f4d41cd614655e700 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
