FROM golang:latest as bot_base_build

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod /app
COPY go.sum /app

RUN go mod download

COPY . .
RUN go build -v -o bot cmd/bot/main.go

FROM alpine:3.15.0

WORKDIR /app

COPY --from=bot_base_build /app/bot .

ENTRYPOINT ["/app/main"]
