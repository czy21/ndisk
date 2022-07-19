FROM alpine:3.16
RUN apk update && apk add tzdata
#COPY ./web/build/ /app/dist/
COPY ./backend/build/ /app/

RUN chmod +x /app/ndisk

CMD ["/app/ndisk"]