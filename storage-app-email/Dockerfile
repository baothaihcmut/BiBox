FROM alpine:latest

WORKDIR /app/
COPY ./bin/main .
COPY ./templates ./templates
RUN chmod +x /app/main
EXPOSE 8080

CMD ["/app/main"]
