# grab the caddy binary
FROM caddy:2.2.1-alpine AS caddy

# build the api
FROM golang:1.15-alpine AS builder
WORKDIR /go/src/github.com/craftcms/nitro
COPY . .
RUN GOOS=linux go build -o api ./cmd/api

# build the final image
FROM alpine:3.12
RUN apk --no-cache add ca-certificates nss-tools supervisor
COPY --from=caddy /usr/bin/caddy /usr/bin/caddy
COPY --from=builder /go/src/github.com/craftcms/nitro/api /usr/bin/nitrod
COPY .docker/supervisor.conf /etc/supervisor/conf.d/supervisor.conf
ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisor.conf"]
EXPOSE 443 80 5000
