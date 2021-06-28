FROM ubuntu:20.04

ARG IMAGE_USER
ENV IMAGE_USER=${IMAGE_USER}

ARG DEFAULT_PHP_VERSION=8.0
ENV DEFAULT_PHP_VERSION=${DEFAULT_PHP_VERSION}

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y \
    supervisor \
    software-properties-common && \
    add-apt-repository ppa:ondrej/php && \
    add-apt-repository ppa:nginx/stable && \
    apt-get update && \
    apt-get remove -y software-properties-common && \
    apt-get install -y \
    php${DEFAULT_PHP_VERSION}-bcmath \
    php${DEFAULT_PHP_VERSION}-cli \
    php${DEFAULT_PHP_VERSION}-curl \
    php${DEFAULT_PHP_VERSION}-gd \
    php${DEFAULT_PHP_VERSION}-fpm \
    php${DEFAULT_PHP_VERSION}-iconv \
    php${DEFAULT_PHP_VERSION}-imagick \
    php${DEFAULT_PHP_VERSION}-intl \
    php${DEFAULT_PHP_VERSION}-mbstring \
    php${DEFAULT_PHP_VERSION}-mysql \
    php${DEFAULT_PHP_VERSION}-opcache \
    php${DEFAULT_PHP_VERSION}-pgsql \
    php${DEFAULT_PHP_VERSION}-redis \
    php${DEFAULT_PHP_VERSION}-soap \
    php${DEFAULT_PHP_VERSION}-xdebug \
    php${DEFAULT_PHP_VERSION}-zip \
    php${DEFAULT_PHP_VERSION}-fpm \
    jpegoptim \
    nginx && \
    apt upgrade -y && \
    apt-get autoremove -y

# install database tools
RUN apt-get install -y --no-install-recommends mysql-client postgresql-client

# create the user and home directory
RUN mkdir /app && \
    addgroup ${IMAGE_USER} --gid 1000 && \
    adduser ${IMAGE_USER} --home /app --uid 1000 --gid 1000 --no-create-home && \
    chown -R ${IMAGE_USER}:${IMAGE_USER} /app

WORKDIR /app

# supervisor config

# copy the config
COPY supervisor.conf /etc/supervisor/conf.d/supervisor.conf
# create the pid for supervisor and set user permissions
RUN touch /run/supervisord.pid && chown ${IMAGE_USER}:${IMAGE_USER} /run/supervisord.pid

# nginx config

# copy the default config and includes
COPY nitro/general.conf /etc/nginx/nitro/
COPY nitro/php_fastcgi.conf /etc/nginx/nitro/
COPY nitro/security.conf /etc/nginx/nitro/
COPY nginx/zz_default.conf /etc/nginx/sites-available/default
# comment out the nginx user from the config
RUN sed -i "s/user www-data;.*/#user www-data;/g" /etc/nginx/nginx.conf
# create the pid for the nginx config
RUN touch /var/run/nginx.pid && \
    chown -R ${IMAGE_USER}:${IMAGE_USER} /var/run/nginx.pid
# make the /var/lib/nginx/body directory and set permissions
RUN mkdir -p /var/lib/nginx/body && chown -R ${IMAGE_USER}:${IMAGE_USER} /var/lib/nginx
# set the permissions for the user on nginx directories
RUN chown -R ${IMAGE_USER}:${IMAGE_USER} /etc/nginx/conf.d/ && \
    chown -R ${IMAGE_USER}:${IMAGE_USER} /var/log/nginx

# php-fpm config

# set permissions for the user on the php logs
RUN touch /var/log/php${DEFAULT_PHP_VERSION}-fpm.log && chown -R ${IMAGE_USER}:${IMAGE_USER} /var/log/php${DEFAULT_PHP_VERSION}-fpm.log
# set permissions on php-fpm directories
RUN mkdir -p /run/php/ && touch /run/php/php${DEFAULT_PHP_VERSION}-fpm.sock && chown -R ${IMAGE_USER}:${IMAGE_USER} /run/php/php${DEFAULT_PHP_VERSION}-fpm.sock
RUN touch /run/php/php${DEFAULT_PHP_VERSION}-fpm.pid && chown ${IMAGE_USER}:${IMAGE_USER} /run/php/php${DEFAULT_PHP_VERSION}-fpm.pid
# use our php-fpm conf
COPY www.conf /etc/php/8.0/fpm/pool.d/www.conf

USER ${IMAGE_USER}

EXPOSE 80 443

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisor.conf"]
