# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:294846130303560fa9cb1fcce26b0fb1c35007906abeac9419dce2fc13a07286 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
