FROM golang:1.14

LABEL maintainer="Ryan Siu <ryansiu1995@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o gcb-visualizer .

ENTRYPOINT ["./gcb-visualizer"]
