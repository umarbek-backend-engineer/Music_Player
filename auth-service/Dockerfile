FROM golang:1.25 as builder

WORKDIR /app

# copy go mod
COPY go.mod go.sum ./
RUN go mod download


#copy all files
COPY . . /app/


#build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# run the stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/.env .
COPY migrator/migrations /migrations

EXPOSE 50053

CMD [ "./app" ]