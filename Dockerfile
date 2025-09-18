FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o taskscheduler main.go

ENTRYPOINT ["./taskscheduler"]
