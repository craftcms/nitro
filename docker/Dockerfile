FROM caddy/caddy:2-alpine
RUN apk update --no-cache && apk add nss-tools
COPY Caddyfile /etc/caddy/Caddyfile
CMD ["caddy", "run", "--resume", "--config=/etc/caddy/Caddyfile"]
