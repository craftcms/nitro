# build the caddy binary

# build the nitro gRPC API


FROM caddy/caddy:2-alpine
COPY Caddyfile /etc/caddy/Caddyfile
RUN mkdir --parents /var/www/html
COPY index.html /var/www/html/
CMD ["caddy", "run", "--resume", "--config=/etc/caddy/Caddyfile"]

# build the final image
# FROM alpine:3.12
# RUN apk update --no-cache && apk add nss-tools ca-certificates
# EXPOSE "80"
# EXPOSE "443"
# EXPOSE "5000"
# ENTRYPOINT [ "/usr/bin/caddy" ]
