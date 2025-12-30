# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:b6ba5234ed128185b0b81070813c77dd5c973aec7703f4646805f11377627408 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
