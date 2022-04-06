FROM golang:1.17-alpine AS build

WORKDIR /app


COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main

ENTRYPOINT [ "./main" ]    