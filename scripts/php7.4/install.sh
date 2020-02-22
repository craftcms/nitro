#!/bin/bash
# add ppas
sudo add-apt-repository -y ppa:nginx/development
sudo add-apt-repository -y ppa:ondrej/php
sudo apt-get update

# install nginx and php7.4
sudo apt-get install -y nginx php7.4 php7.4-mbstring php7.4-cli php7.4-curl php7.4-fpm php7.4-gd php7.4-intl php7.4-json php7.4-mysql \
    php7.4-opcache php7.4-pgsql php7.4-zip php7.4-xml redis

# install composer
php -r "readfile('http://getcomposer.org/installer');" | sudo php -- --install-dir=/usr/bin/ --filename=composer
