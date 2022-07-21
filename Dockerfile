FROM golang:1.18.4-bullseye

#COPY ./web/build/ /app/dist/
COPY ./backend/build/ /app/

RUN chmod +x /app/ndisk

CMD ["/app/ndisk"]