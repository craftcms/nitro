FROM debian:bullseye-slim

ARG IMAGE_USER
ENV IMAGE_USER=${IMAGE_USER}

ARG PHP_VERSION=8.0
ENV PHP_VERSION=${PHP_VERSION}

ARG NVM_VERSION=0.38.0
ENV NVM_VERSION=${NVM_VERSION}

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y apt-transport-https build-essential ca-certificates curl git lsb-release nano nodejs supervisor software-properties-common unzip vim wget nginx && \
    wget -O /etc/apt/trusted.gpg.d/php.gpg https://packages.sury.org/php/apt.gpg && \
    sh -c 'echo "deb https://packages.sury.org/php/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/php.list' && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && \
    apt-get remove -y software-properties-common && \
    apt-get install -y \
    php${PHP_VERSION}-bcmath \
    php${PHP_VERSION}-cli \
    php${PHP_VERSION}-curl \
    php${PHP_VERSION}-dom \
    php${PHP_VERSION}-gd \
    php${PHP_VERSION}-fpm \
    php${PHP_VERSION}-iconv \
    php${PHP_VERSION}-imagick \
    php${PHP_VERSION}-intl \
    php${PHP_VERSION}-mbstring \
    php${PHP_VERSION}-mysql \
    php${PHP_VERSION}-opcache \
    php${PHP_VERSION}-pgsql \
    php${PHP_VERSION}-redis \
    php${PHP_VERSION}-soap \
    php${PHP_VERSION}-xdebug \
    php${PHP_VERSION}-zip \
    php${PHP_VERSION}-fpm \
    # image processing utils
    gifsicle \
    libjpeg-turbo-progs \
    libwebp6 \
    jpegoptim \
    optipng \
    pngquant && \
    apt-get upgrade -y && \
    apt-get autoremove -y

# install database tools
RUN apt-get install -y --no-install-recommends mariadb-client postgresql-client

# install composer
RUN cd /tmp && curl -sS https://getcomposer.org/installer -o composer-setup.php && \
    HASH=`curl -sS https://composer.github.io/installer.sig` && \
    php -r "if (hash_file('SHA384', 'composer-setup.php') === '$HASH') { echo 'Installer verified'; } else { echo 'Installer corrupt'; unlink('composer-setup.php'); } echo PHP_EOL;" && \
    php composer-setup.php --install-dir=/usr/local/bin --filename=composer

# copy the blackfire script
COPY scripts/blackfire-install.sh /tmp/blackfire-install.sh

# run the installation script
RUN set -ex && sh /tmp/blackfire-install.sh && rm -rf /tmp/blackfire*

# create the user and home directory
RUN mkdir /app && \
    addgroup ${IMAGE_USER} --gid 1000 && \
    adduser ${IMAGE_USER} --uid 1000 --gid 1000 && \
    chown -R ${IMAGE_USER}:${IMAGE_USER} /app

WORKDIR /app

# install nvm
RUN export XDG_CONFIG_HOME=/home/${IMAGE_USER} && \
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v${NVM_VERSION}/install.sh | bash && \
    echo ". /home/${IMAGE_USER}/nvm/nvm.sh" >> /home/${IMAGE_USER}/.bashrc && \
    chown -R ${IMAGE_USER}:${IMAGE_USER} /home/${IMAGE_USER}

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
RUN rm -rf /etc/nginx/sites-available/default && rm -rf /etc/nginx/sites-enabled/default
COPY nginx/zz_default.conf /etc/nginx/sites-available/zz_default
RUN ln -s /etc/nginx/sites-available/zz_default /etc/nginx/sites-enabled/
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
RUN touch /var/log/php${PHP_VERSION}-fpm.log && chown -R ${IMAGE_USER}:${IMAGE_USER} /var/log/php${PHP_VERSION}-fpm.log
# set permissions on php-fpm directories
RUN mkdir -p /run/php/ && touch /run/php/php${PHP_VERSION}-fpm.sock && chown -R ${IMAGE_USER}:${IMAGE_USER} /run/php/php${PHP_VERSION}-fpm.sock
RUN touch /run/php/php${PHP_VERSION}-fpm.pid && chown ${IMAGE_USER}:${IMAGE_USER} /run/php/php${PHP_VERSION}-fpm.pid
# use our php-fpm conf
COPY www.conf /etc/php/${PHP_VERSION}/fpm/pool.d/www.conf
# copy custom.ini settings
COPY php/nitro.ini /etc/php/${PHP_VERSION}/cli/conf.d
COPY php/nitro.ini /etc/php/${PHP_VERSION}/fpm/conf.d

USER ${IMAGE_USER}

EXPOSE 80 443 3000-3005

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisor.conf"]
