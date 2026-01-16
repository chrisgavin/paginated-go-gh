# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:bc45dfd319e982dffe4de14428c77defe5b938e29d9bc6edfbc0b9a1fc171cb3 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
