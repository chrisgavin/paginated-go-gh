# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:d52df9c279840adf958d017ebb275651ed8338b953d39817bc3633a2e6b1bbcc AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
