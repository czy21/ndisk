FROM golang:1.18.4-alpine

#COPY ./web/build/ /app/dist/
COPY ./backend/build/ /app/

RUN chmod +x /app/ndisk

CMD ["/app/ndisk"]