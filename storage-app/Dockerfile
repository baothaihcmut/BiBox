FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/storage-app ./...

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/storage-app ./main

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]
