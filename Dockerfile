ARG GO_VERSION=1.16

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app

# download and cache deps
COPY go.mod .
COPY go.sum .
RUN go mod download

# build the app
COPY . .
RUN go build -o ./app ./server.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /app .

# expose HTTP server port
EXPOSE 1323

# start the server
ENTRYPOINT ["./app"]