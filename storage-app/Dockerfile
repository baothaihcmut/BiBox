FROM alpine:latest

WORKDIR /app/
COPY ./bin/storage-app main
RUN chmod +x /app/main
EXPOSE 8080

CMD ["./main"]
