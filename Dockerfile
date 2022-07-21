FROM buildpack-deps:bullseye-scm

COPY ./backend/build/ /app/

RUN chmod +x /app/ndisk

CMD ["/app/ndisk"]
