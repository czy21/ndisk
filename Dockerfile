FROM debian:bullseye-slim

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		libpython3-stdlib \
	; \
	rm -rf /var/lib/apt/lists/*

COPY ./backend/build/ /app/

RUN chmod +x /app/ndisk

CMD ["/app/ndisk"]
