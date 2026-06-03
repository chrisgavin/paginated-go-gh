# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:68cb6d68bed024785b69195b89af7ac7a444f27791435f98647edff595aa0479 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
