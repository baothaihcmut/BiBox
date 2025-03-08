# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}
RUN go build -o main ./cmd/storage-app/main.go
FROM alpine:latest

WORKDIR /app/
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
