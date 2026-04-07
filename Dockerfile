# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:cd78d88e00afadbedd272f977d375a6247455f3a4b1178f8ae8bbcb201743a8a AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
