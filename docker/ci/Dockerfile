FROM golang:alpine

RUN apk update && \
  apk upgrade --no-cache && \
  apk add --no-cache \
    build-base \
    curl \
    git \
    python3 \
    nodejs \
    yarn

RUN curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh

RUN go get github.com/markbates/pkger/cmd/pkger