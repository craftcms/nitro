# grab the caddy binary
FROM caddy:2.2.1-alpine AS caddy

# build the api
FROM golang:1.15-alpine AS builder
WORKDIR /go/src/github.com/craftcms/nitro
COPY . .
RUN go build -ldflags="-s -w" -o nitrod ./cmd/api

# build the final image
FROM alpine:3.12
RUN apk --no-cache add ca-certificates nss-tools supervisor
RUN mkdir --parents /var/www/html
RUN mkdir --parents /etc/caddy/
COPY .docker/Caddyfile /etc/caddy/Caddyfile
COPY .docker/index.html /var/www/html/
COPY --from=caddy /usr/bin/caddy /usr/bin/caddy
COPY --from=builder /go/src/github.com/craftcms/nitro/nitrod /usr/bin/nitrod

COPY .docker/supervisor.conf /etc/supervisor/conf.d/supervisor.conf

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisor.conf"]
EXPOSE 443 80 5000
# TODO remove this after testing
EXPOSE 2019
