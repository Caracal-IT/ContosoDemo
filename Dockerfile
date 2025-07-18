FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o contoso-server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/contoso-server .
COPY ./public ./public

EXPOSE 8080

CMD ["./contoso-server"]

