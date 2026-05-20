# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:8530a4fd262b294857b6c1d07203f57261863cafe5924f90ebdb4ba96b16ad15 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
