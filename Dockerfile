FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

RUN go build -o contoso_server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/contoso_server .
COPY ./public ./public

EXPOSE 8080

CMD ["./contoso_server"]
