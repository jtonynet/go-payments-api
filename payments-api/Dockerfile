FROM golang:1.23-alpine AS builder
WORKDIR /usr/src/app

# Environment variables which CompileDaemon requires to run
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0


COPY go.mod ./
COPY go.sum ./

COPY . .

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
