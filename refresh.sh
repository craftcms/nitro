#!/usr/bin/env bash
version="$1"

if [ -z "$version" ]; then
  echo "you must specify a version of nitro"
  exit 1
fi

# always download and install the nitrod api regardless of the version passed
echo "installing nitrod and nitrod.service"
curl -s https://api.github.com/repos/craftcms/nitro/releases/latest \
  | grep "browser_download_url" \
  | grep "nitrod_linux_x86_64" \
  | cut -d : -f 2,3 | tr -d \" \
  | wget --directory-prefix=/tmp -qi -
cd /tmp && tar xfz /tmp/nitrod_linux_x86_64.tar.gz
mv /tmp/nitrod /usr/sbin/
mv /tmp/nitrod.service /etc/systemd/system/
systemctl daemon-reload
systemctl start nitrod
systemctl enable nitrod

if [ "$version" == "1.0.0-RC1" ] || [ "$version" == "1.0.0-RC1.1" ]; then
  echo "running script for 1.0.0-RC1 and 1.0.0-RC1.1"

  echo "updating /etc/resolv.conf"
  sed -i 's|nameserver 1.1.1.1|nameserver 127.0.0.53\nnameserver 1.1.1.1\nnameserver 1.0.0.1\nnameserver 8.8.8.8\nnameserver 8.8.4.4|g' /etc/resolv.conf

  echo "installing nitrod and nitrod.service"
  curl -s https://api.github.com/repos/craftcms/nitro/releases/latest \
    | grep "browser_download_url" \
    | grep "nitrod_linux_x86_64" \
    | cut -d : -f 2,3 | tr -d \" \
    | wget --directory-prefix=/tmp -qi -
  cd /tmp && tar xfz /tmp/nitrod_linux_x86_64.tar.gz
  mv /tmp/nitrod /usr/sbin/
  mv /tmp/nitrod.service /etc/systemd/system/
  systemctl daemon-reload
  systemctl start nitrod
  systemctl enable nitrod
fi

# script for beta 7
if [ "$version" == "1.0.0-beta.7" ] || [ "$version" == "1.0.0-beta.8" ] || [ "$version" == "1.0.0-beta.9" ] || [ "$version" == "1.0.0-beta.10" ]; then
  echo "running script for 1.0.0-beta.9"

  cat >"/opt/nitro/nginx/template.conf" <<-EndOfMessage
# Hat tip to https://github.com/nystudio107/nginx-craft

server {
  # Listen for both IPv4 & IPv6 on port 80
  listen 80;
  listen [::]:80;

  # General virtual host settings
  server_name CHANGESERVERNAME;
  root CHANGEWEBROOTDIR;
  index index.html index.htm index.php;
  charset utf-8;

  # Enable serving of static gzip files as per: http://nginx.org/en/docs/http/ngx_http_gzip_static_module.html
  gzip_static  on;

  # Enable server-side includes as per: http://nginx.org/en/docs/http/ngx_http_ssi_module.html
  ssi on;

  # Disable limits on the maximum allowed size of the client request body
  client_max_body_size 0;

  # 404 error handler
  error_page 404 /index.php\$is_args\$args;

  # Root directory location handler
  location / {
    try_files \$uri/index.html \$uri \$uri/ /index.php\$is_args\$args;
  }

  # php-fpm configuration
  location ~ [^/]\.php(/|$) {
    include snippets/fastcgi-php.conf;

    fastcgi_pass unix:/var/run/php/phpCHANGEPHPVERSION-fpm.sock;

    # FastCGI params
    fastcgi_param CRAFT_NITRO 1;
    fastcgi_param DB_USER nitro;
    fastcgi_param DB_PASSWORD nitro;
    fastcgi_param HTTP_PROXY "";
    fastcgi_param HTTP_HOST CHANGESERVERNAME;

    # Don't allow browser caching of dynamically generated content
    add_header Last-Modified \$date_gmt;
    add_header Cache-Control "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0";
    if_modified_since off;
    expires off;
    etag off;

    fastcgi_intercept_errors off;
    fastcgi_buffer_size 16k;
    fastcgi_buffers 4 16k;
    fastcgi_connect_timeout 240;
    fastcgi_send_timeout 240;
    fastcgi_read_timeout 240;
  }

  # Disable reading of Apache .htaccess files
  location ~ /\.ht {
    deny all;
  }

  # Misc settings
  sendfile off;
}
EndOfMessage
fi

# scripts for beta 5
if [ "$version" == "1.0.0-beta.5" ] || [ "$version" == "1.0.0-beta.6" ]; then
  echo "running sync script for 1.0.0-beta.6"

  echo "installing mysql-client and postgresql-client tools"
  apt install -y mysql-client postgresql-client

  echo "removing the NGINX ppa"
  add-apt-repository --remove ppa:nginx/stable
  apt remove nginx
  apt update
  apt upgrade -y
  apt install -y nginx

  echo "setting the default mysql conf for 5.x"
  mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d/5/
  cat >"/home/ubuntu/.nitro/databases/mysql/conf.d/5/mysql.conf" <<-EndOfMessage
[mysqld]
max_allowed_packet=256M
wait_timeout=86400
default-authentication-plugin=mysql_native_password
EndOfMessage

  echo "setting the default mysql conf for 8.x"
  mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d/8/
  cat >"/home/ubuntu/.nitro/databases/mysql/conf.d/8/mysql.conf" <<-EndOfMessage
[mysqld]
max_allowed_packet=256M
wait_timeout=86400
default-authentication-plugin=mysql_native_password
[mysqldump]
column-statistics=0
EndOfMessage

  echo "refresh script has completed!"
  exit 0
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
  chown -R ubuntu:ubuntu /home/ubuntu/.nitro

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
    error_page 404 /index.php?\$query_string;
    ssi on;
    server_name CHANGESERVERNAME;
    client_max_body_size 100M;
    location / {
      try_files \$uri \$uri/ /index.php\$is_args\$args;
    }
    location ~ \.php\$ {
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
