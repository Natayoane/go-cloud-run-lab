FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=build /app/cloudrun .
COPY .env .env

EXPOSE 8080

ENTRYPOINT ["./cloudrun"]