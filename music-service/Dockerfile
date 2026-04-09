# Build stage
FROM golang:1.25 AS builder

WORKDIR /app


# Copy go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy all files
COPY . .


# Build binary
RUN go build -o app ./cmd

# Run stage
FROM alpine:latest
    

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 50051

CMD ["./app"]