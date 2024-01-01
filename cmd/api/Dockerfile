FROM golang:latest as api_base_build

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod /app
COPY go.sum /app

RUN go mod download

COPY . .
RUN go build -v -o api cmd/api/main.go

FROM alpine:3.15.0

WORKDIR /app

COPY --from=api_base_build /app/api .

EXPOSE 8001

ENTRYPOINT ["/app/main"]