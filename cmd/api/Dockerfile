FROM golang:latest

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod /app
COPY go.sum /app

RUN go mod download

COPY . .
RUN go build -o /app/bin/api cmd/api/main.go

EXPOSE 8001

ENTRYPOINT ["/app/bin/api"]