FROM caddy/caddy:2-alpine
RUN apk update --no-cache && apk add nss-tools ca-certificates
COPY Caddyfile /etc/caddy/Caddyfile
RUN mkdir --parents /var/www/html
COPY index.html /var/www/html/
CMD ["caddy", "run", "--resume", "--config=/etc/caddy/Caddyfile"]
