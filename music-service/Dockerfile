# Build stage
FROM golang:1.25 AS builder

WORKDIR /app


# Copy go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy all files
COPY . .


# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# Run stage
FROM alpine:latest
    

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/.env .
COPY /migrations /migrations

EXPOSE 50051

CMD ["./app"]