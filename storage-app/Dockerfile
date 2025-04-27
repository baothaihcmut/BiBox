FROM golang:1.23 AS builder

WORKDIR /app

COPY libs/go.mod libs/go.sum ./libs/

COPY storage-app/go.mod storage-app/go.sum ./storage-app/


RUN go work init ./storage-app ./libs

RUN go mod download -C ./storage-app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/storage-app storage-app/cmd/storage-app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/storage-app ./main

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]
