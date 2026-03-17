# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:16e774b791968123d6af5ba4eec19cf91c4208cb1f5849efda5d4ffaf6d1c038 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
