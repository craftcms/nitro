package command

var cloudInit = `#cloud-config
packages:
  - redis
  - jq
write_files:
  # create the main bootstrap script
  - path: /opt/nitro/bootstrap.sh
    content: |
      #!/usr/bin/env bash
      phpversion="$1"
      database="$2"

      # run the correct script depending on the version of PHP
      if [ "$phpversion" == '7.3' ]; then
        sudo echo "7.3" > /opt/nitro/php_version
        sudo bash /opt/nitro/php/php-73.sh
      elif [ "$phpversion" == '7.2' ]; then
        sudo echo "7.2" > /opt/nitro/php_version
        sudo bash /opt/nitro/php/php-72.sh
      else
        sudo echo "7.4" > /opt/nitro/php_version
        sudo bash /opt/nitro/php/php-74.sh
      fi

      # run the correct script depending on the database
      if [ "$database" == 'postgres' ]; then
        # install postgres
        sudo bash /opt/nitro/postgres/install.sh
        # set the engine
        sudo echo "postgres" > /opt/nitro/db_engine
        # run the postgres setup
        sudo bash /opt/nitro/postgres/setup.sh
      else
        # install mariadb
        sudo bash /opt/nitro/mariadb/install.sh
        # set the engine
        sudo echo "mariadb" > /opt/nitro/db_engine
        # run the mariadb setup
        sudo bash /opt/nitro/mariadb/setup.sh
      fi

      # install nginx
      sudo bash /opt/nitro/nginx/install.sh

      # remove the default site
      if [ -f /etc/nginx/sites-enabled/default ]; then
          sudo rm /etc/nginx/sites-enabled/default
      fi

      # change user php is running as
      export phpversion=$(cat /opt/nitro/php_version)
      sudo sed -i "s|user = www-data|user = ubuntu|g" /etc/php/"$phpversion"/fpm/pool.d/www.conf
      sudo sed -i "s|group = www-data|group = ubuntu|g" /etc/php/"$phpversion"/fpm/pool.d/www.conf

      # set php.ini settings for craft
      sudo sed -i "s|memory_limit = 128M|memory_limit = 256M|g" /etc/php/"$phpversion"/fpm/php.ini
      sudo sed -i "s|max_execution_time = 30|max_execution_time = 120|g" /etc/php/"$phpversion"/fpm/php.ini

      # set xDebug settings whether it's enabled or not
      sudo sed -i -e "\$axdebug.remote_enable=1\nxdebug.remote_connect_back=0\nxdebug.remote_host=localhost\nxdebug.remote_port=9000\nxdebug.remote_log=/var/log/nginx/xdebug.log" /etc/php/"$phpversion"/mods-available/xdebug.ini

      sudo service php"$phpversion"-fpm restart
  - path: /opt/nitro/refresh.sh
    content: |
      #!/usr/bin/env bash
      export FILE_CONTENT="$1"
      export FILE_PATH="$2"
      if [ -n "$FILE_CONTENT" ]; then
         echo "$FILE_CONTENT" | sudo tee "$FILE_PATH"
      else
         echo "content was empty, skipping"
      fi
  - path: /opt/nitro/update.sh
    content: |
      #!/usr/bin/env bash
      sudo apt update -y && sudo apt upgrade -y
  - path: /opt/nitro/php/enable-xdebug.sh
    content: |
      #!/bin/bash
      export phpversion=$(cat /opt/nitro/php_version)
      sudo phpenmod -v "$phpverison" xdebug
      echo "enabled xdebug for $phpversion"
  - path: /opt/nitro/php/disable-xdebug.sh
    content: |
      #!/bin/bash
      export phpversion=$(cat /opt/nitro/php_version)
      sudo phpdismod -v "$phpverison" xdebug
      echo "disabled xdebug for $phpversion"
  # create the php install scripts
  - path: /opt/nitro/php/php-74.sh
    content: |
      #!/bin/bash
      apt install -y php7.4 php7.4-mbstring php7.4-cli php7.4-curl php7.4-fpm php7.4-gd php7.4-intl php7.4-json \
      php7.4-mysql php7.4-opcache php7.4-pgsql php7.4-zip php7.4-xml php-xdebug php-imagick
  - path: /opt/nitro/php/php-73.sh
    content: |
      #!/bin/bash
      apt install -y php7.3 php7.3-mbstring php7.3-cli php7.3-curl php7.3-fpm php7.3-gd php7.3-intl php7.3-json \
      php7.3-mysql php7.3-opcache php7.3-pgsql php7.3-zip php7.3-xml php-xdebug php-imagick
  - path: /opt/nitro/php/php-72.sh
    content: |
      #!/bin/bash
      apt install -y php7.2 php7.2-mbstring php7.2-cli php7.2-curl php7.2-fpm php7.2-gd php7.2-intl php7.2-json \
      php7.2-mysql php7.2-opcache php7.2-pgsql php7.2-zip php7.2-xml php-xdebug php-imagick
  # create the composer install script
  - path: /opt/nitro/composer-install.sh
    content: |
      # install composer
      php -r "readfile('http://getcomposer.org/installer');" | sudo php -- --install-dir=/usr/bin/ --filename=composer
  # create mariadb scripts
  - path: /opt/nitro/mariadb/setup.sh
    content: |
      # remove bind from mariadb to allow remote connection
      systemctl stop mariadb
      sed -i 's|bind-address|#bind-address|g' /etc/mysql/mariadb.conf.d/50-server.cnf
      systemctl start mariadb

      # create the database and user
      sudo mysql -u root -e "CREATE DATABASE IF NOT EXISTS craftcms;"
      sudo mysql -u root -e "CREATE USER IF NOT EXISTS 'nitro'@'%' IDENTIFIED BY 'nitro';"
      sudo mysql -u root -e "GRANT CREATE, ALTER, INDEX, LOCK TABLES, REFERENCES, UPDATE, DELETE, DROP, SELECT, INSERT ON *.* TO 'nitro'@'%';"
      sudo mysql -u root -e "FLUSH PRIVILEGES;"
  - path: /opt/nitro/mariadb/install.sh
    content: |
      #!/bin/bash
      apt install -y mariadb-server
  - path: /opt/nitro/postgres/setup.sh
    content: |
      # allow remote access to postgres
      sudo sed -i "s|#listen_addresses = 'localhost'|listen_addresses = '\*'|g" /etc/postgresql/10/main/postgresql.conf
      sudo sed -i 's|127.0.0.1/32|0.0.0.0/0|g' /etc/postgresql/10/main/pg_hba.conf
      sudo service postgresql restart

      # create the user and database
      sudo su - postgres -c "createuser --createdb --login nitro"
      sudo -u postgres psql -c "ALTER USER nitro WITH PASSWORD 'nitro';"
      sudo su - postgres -c "createdb craftcms"
      sudo -u postgres psql -c "GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO nitro;"
  - path: /opt/nitro/postgres/install.sh
    content: |
      #!/bin/bash
      apt install -y postgresql postgresql-contrib
  # create nginx install scripts
  - path: /opt/nitro/nginx/install.sh
    content: |
      #!/bin/bash
      apt install -y nginx
  - path: /opt/nitro/nginx/template.conf
    content: |
      server {
          listen 80;
          listen [::]:80;

          root /home/ubuntu/sites/CHANGEPATH/CHANGEPUBLICDIR;

          index index.php;
          gzip_static  on;
          error_page 404 /index.php?$query_string;
          ssi on;
          server_name CHANGESERVERNAME.test;

          location / {
              try_files $uri $uri/ /index.php$is_args$args;
          }

          location ~ \.php$ {
             include snippets/fastcgi-php.conf;
             fastcgi_pass unix:/var/run/php/phpCHANGEPHPVERSION-fpm.sock;
          }
      }
  - path: /opt/nitro/nginx/add-site.sh
    content: |
      #!/usr/bin/env bash
      NEW_HOST_NAME="$1"
      REQUESTED_PHP_VERSION="$2"
      PUBLIC_DIR="$3"
      export INSTALLED_PHP_VERSION=$(cat /opt/nitro/php_version)

      if [ "$REQUESTED_PHP_VERSION" != "$INSTALLED_PHP_VERSION" ]; then
          if [ "$REQUESTED_PHP_VERSION" == '7.3' ]; then
              export REQUESTED_PHP_VERSION="7.3"
              sudo bash /opt/nitro/php/php-73.sh
          elif [ "$REQUESTED_PHP_VERSION" == '7.2' ]; then
              export REQUESTED_PHP_VERSION="7.2"
              sudo bash /opt/nitro/php/php-72.sh
          else
              export REQUESTED_PHP_VERSION="7.4"
              sudo bash /opt/nitro/php/php-74.sh
          fi
      fi

      # make the directories
      mkdir -p /home/ubuntu/sites/"$NEW_HOST_NAME"

      # copy the nitro nginx template into sites-available/default
      sudo cp /opt/nitro/nginx/template.conf /etc/nginx/sites-available/"$NEW_HOST_NAME"

      # change the variables
      sudo sed -i 's|CHANGEPATH|'$NEW_HOST_NAME'|g' /etc/nginx/sites-available/"$NEW_HOST_NAME"
      sudo sed -i 's|CHANGESERVERNAME|'$NEW_HOST_NAME'|g' /etc/nginx/sites-available/"$NEW_HOST_NAME"
      sudo sed -i 's|CHANGEPUBLICDIR|'$PUBLIC_DIR'|g' /etc/nginx/sites-available/"$NEW_HOST_NAME"
      sudo sed -i 's|CHANGEPHPVERSION|'$REQUESTED_PHP_VERSION'|g' /etc/nginx/sites-available/"$NEW_HOST_NAME"

      # enable the nginx site
      sudo ln -s /etc/nginx/sites-available/"$NEW_HOST_NAME" /etc/nginx/sites-enabled/

      # reload nginx
      sudo service nginx reload
  - path: /opt/nitro/nginx/remove-site.sh
    content: |
      #!/usr/bin/env bash
      REMOVE_HOST="$1"
      # remove the nginx site
      sudo rm /etc/nginx/sites-enabled/"$REMOVE_HOST"
      # reload nginx
      sudo service nginx reload
  - path: /opt/nitro/nginx/tail-logs.sh
    content: |
      #!/usr/bin/env bash
      sudo tail -f /var/log/nginx/access.log -f /var/log/nginx/error.log
runcmd:
  - sudo add-apt-repository -y ppa:nginx/stable
  - sudo add-apt-repository -y ppa:ondrej/php
  - sudo apt update -y
`
