#!/usr/bin/env bash
version="$1"

if [ -z "$version" ]; then
  echo "you must specify a version of nitro"
  exit 1
fi

if [ "$version" == "1.0.0-beta.5" ]; then
  echo "running sync script for 1.0.0-beta.5"

  echo "setting the default mysql conf"
  cat >"/home/ubuntu/.nitro/databases/mysql/conf.d/mysql.conf" <<-EndOfMessage
[mysqld]
max_allowed_packet=256M
wait_timeout=86400
default-authentication-plugin=mysql_native_password
EndOfMessage
fi

# beta 3 and beta 4 scripts
if [ "$version" == "1.0.0-beta.3" ] || [ "$version" == "1.0.0-beta.4" ]; then
  echo "running sync script for 1.0.0-beta.3"

  echo "ensuring composer is installed"
  export COMPOSER_HOME=/home/ubuntu/composer
  curl -sS https://getcomposer.org/installer -o composer-setup.php
  php composer-setup.php --install-dir=/usr/local/bin --filename=composer
  rm composer-setup.php

  # copy skeleton
  cp /etc/skel/.bashrc /home/ubuntu/.bashrc
  cp /etc/skel/.profile /home/ubuntu/.profile
  cp /etc/skel/.bash_logout /home/ubuntu/.bash_logout

  # create setup scripts
  mkdir -p /home/ubuntu/sites
  mkdir -p /home/ubuntu/.nitro/databases/imports
  mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d
  mkdir -p /home/ubuntu/.nitro/databases/mysql/backups
  mkdir -p /home/ubuntu/.nitro/databases/postgres/conf.d
  mkdir -p /home/ubuntu/.nitro/databases/postgres/backups
  chown -R ubuntu:ubuntu /home/ubuntu/

  # create the files
  echo "setting the default mysql conf"
  cat >"/home/ubuntu/.nitro/databases/mysql/conf.d/mysql.conf" <<-EndOfMessage
[mysqld]
max_allowed_packet=256M
wait_timeout=86400
default-authentication-plugin=mysql_native_password
EndOfMessage

  echo "setting the default mysql setup"
  cat >"/home/ubuntu/.nitro/databases/mysql/setup.sql" <<-EndOfMessage
CREATE USER IF NOT EXISTS 'nitro'@'localhost' IDENTIFIED BY 'nitro';
CREATE USER IF NOT EXISTS 'nitro'@'%' IDENTIFIED BY 'nitro';
GRANT ALL PRIVILEGES ON *.* TO 'nitro'@'localhost' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'nitro'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
EndOfMessage

  echo "setting the default postgres setup"
  cat >"/home/ubuntu/.nitro/databases/postgres/setup.sql" <<-EndOfMessage
ALTER USER nitro WITH SUPERUSER;
EndOfMessage

  echo "updating the nginx template"
  cat >"/opt/nitro/nginx/template.conf" <<-EndOfMessage
  server {
    listen 80;
    listen [::]:80;
    root CHANGEWEBROOTDIR;
    index index.php;
    gzip_static  on;
    error_page 404 /index.php?$query_string;
    ssi on;
    server_name CHANGESERVERNAME;
    client_max_body_size 100M;
    location / {
      try_files $uri $uri/ /index.php$is_args$args;
    }
    location ~ \.php$ {
      include snippets/fastcgi-php.conf;
      fastcgi_pass unix:/var/run/php/phpCHANGEPHPVERSION-fpm.sock;
      fastcgi_read_timeout 240;
      fastcgi_param CRAFT_NITRO 1;
      fastcgi_param DB_USER nitro;
      fastcgi_param DB_PASSWORD nitro;
    }
}
EndOfMessage

  echo "setting DB_USER and DB_PASSWORD environment variables"
  echo "CRAFT_NITRO=1" >>"/etc/environment"
  echo "DB_USER=nitro" >>"/etc/environment"
  echo "DB_PASSWORD=nitro" >>"/etc/environment"

  echo "removing old scripts"
  rm -rf /opt/nitro/scripts

  echo "setup script has completed!"
  exit 0
fi
