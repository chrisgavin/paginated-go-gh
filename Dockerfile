# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:bccbee38ded684484284fdbb72c7d9bc93992f38a0f3b1615a005aa312fa9c4b AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...
