FROM golang:latest

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod /app
COPY go.sum /app

RUN go mod download

COPY . .
RUN go build -v -o /app/bin/bot cmd/bot/main.go

ENTRYPOINT ["/app/bin/bot"]
