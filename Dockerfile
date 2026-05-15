FROM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git build-base

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o mysql-dump-backup .

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates mysql-client tzdata

COPY --from=builder /app/mysql-dump-backup .

CMD ["./mysql-dump-backup"]