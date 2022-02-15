# grab the caddy binary
FROM caddy:2.3.0-alpine AS caddy

# build the api
FROM golang:1.17-alpine AS builder
ARG NITRO_VERSION
ENV NITRO_VERSION=${NITRO_VERSION}
WORKDIR /go/src/github.com/craftcms/nitro
COPY . .
RUN GOOS=linux go build -ldflags="-s -w -X 'github.com/craftcms/nitro/pkg/api/api.Version=${NITRO_VERSION}'" -o nitrod ./cmd/nitrod

# build the final image
FROM alpine:3

# See https://caddyserver.com/docs/conventions#file-locations for details
ENV XDG_CONFIG_HOME /config
ENV XDG_DATA_HOME /data

# label the container
LABEL org.opencontainers.image.title="Craft Nitro"
LABEL org.opencontainers.image.description="Nitro is a command-line tool focused on making local Craft CMS development quick and easy"
LABEL org.opencontainers.image.url=https://getnitro.sh
LABEL org.opencontainers.image.documentation=https://craftcms.com/docs/nitro
LABEL org.opencontainers.image.vendor="Craft CMS"
LABEL org.opencontainers.image.source="https://github.com/craftcms/nitro"

RUN apk --no-cache add ca-certificates nss-tools supervisor postgresql-client mysql-client jq
RUN mkdir --parents /var/www/html
RUN mkdir --parents /etc/caddy/
RUN mkdir --parents /config
RUN mkdir --parents /data

COPY .docker/Caddyfile /etc/caddy/Caddyfile
COPY .docker/index.html /var/www/html/
COPY --from=caddy /usr/bin/caddy /usr/bin/caddy
COPY --from=builder /go/src/github.com/craftcms/nitro/nitrod /usr/bin/nitrod
COPY .docker/supervisor.conf /etc/supervisor/conf.d/supervisor.conf

VOLUME /config
VOLUME /data

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisor.conf"]

EXPOSE 443 80 5000 3000 3001
