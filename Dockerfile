# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:75ca321c953ce0572c709fa186cff872a510a22d6dad515291c6eb29edb9c849 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
