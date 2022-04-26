#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN go build -v -o /go/bin/app/ ./...
