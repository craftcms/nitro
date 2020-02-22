#!/bin/bash
# add ppas
sudo add-apt-repository -y ppa:nginx/development
sudo add-apt-repository -y ppa:ondrej/php
sudo apt-get update

# install nginx, php7.4, and mariadb
sudo apt-get install -y nginx php7.4 php7.4-mbstring php7.4-cli php7.4-curl php7.4-fpm php7.4-gd php7.4-intl php7.4-json php7.4-mysql \
    php7.4-opcache php7.4-pgsql php7.4-zip php7.4-xml redis mariadb-server

# install composer
php -r "readfile('http://getcomposer.org/installer');" | sudo php -- --install-dir=/usr/bin/ --filename=composer

sudo mkdir /etc/dev

# set the passwords
export ROOT_PASS=$(openssl rand -base64 20)
echo "$ROOT_PASS" | sudo tee /etc/dev/.mysql_root_password
export MYSQL_ROOT=$(sudo cat /etc/dev/.mysql_root_password)
export USER_PASS=$(openssl rand -base64 8)
echo "$USER_PASS" | sudo tee /etc/dev/.mysql_user_password
export MYSQL_USER=$(cat /etc/dev/.mysql_user_password)

# run the setup script
sudo sed -i 's|CHANGEME|'$MYSQL_USER'|g' /etc/dev/setup.sql
mysql --user=root --password="$MYSQL_ROOT" < /etc/dev/setup.sql