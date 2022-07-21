FROM alpine:3.16

#COPY ./web/build/ /app/dist/
COPY ./backend/build/ /app/

RUN chmod +x /app/ndisk

CMD ["/app/ndisk"]