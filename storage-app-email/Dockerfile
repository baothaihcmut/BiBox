FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main ./...

FROM alpine:latest

WORKDIR /app/

# Copy binary from the builder
COPY --from=builder /app/bin/main .

# Copy templates folder from builder
COPY --from=builder /app/templates ./templates

# Make binary executable
RUN chmod +x /app/main

EXPOSE 8080

CMD ["/app/main"]
