# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:87a41d2539e5671777734e91f467499ed5eafb1fb1f77221dff2744db7a51775 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
