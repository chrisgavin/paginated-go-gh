# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:b900de91b15b2e2953d930ece1d0ecff0a1590ab2006088d20dcf0f56f1e979f AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
