FROM alpine:3.16

COPY ./web/build/ /app/dist/
COPY ./api/build/ /app/

RUN chmod +x /app/main

CMD ["/app/main"]