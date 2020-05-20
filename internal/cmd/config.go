package cmd

const CloudConfig = `#cloud-config
packages:
  - redis
  - jq
  - apt-transport-https
  - ca-certificates
  - curl
  - gnupg-agent
  - software-properties-common
  - sshfs
  - pv
  - httpie
  - php-cli
  - unzip
write_files:
  - path: /home/ubuntu/.nitro/databases/mysql/conf.d/mysql.cnf
    content: |
      [mysqld]
      max_allowed_packet=256M
      wait_timeout=86400
  - path: /home/ubuntu/.nitro/databases/mysql/setup.sql
    content: |
      CREATE USER IF NOT EXISTS 'nitro'@'localhost' IDENTIFIED BY 'nitro';
      CREATE USER IF NOT EXISTS 'nitro'@'%' IDENTIFIED BY 'nitro';
      GRANT ALL PRIVILEGES ON *.* TO 'nitro'@'localhost' WITH GRANT OPTION;
      GRANT ALL PRIVILEGES ON *.* TO 'nitro'@'%' WITH GRANT OPTION;
      FLUSH PRIVILEGES;
  - path: /home/ubuntu/.nitro/databases/postgres/setup.sql
    content: |
      ALTER USER nitro WITH SUPERUSER;
  - path: /opt/nitro/nginx/template.conf
    content: |
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
  - path: /opt/nitro/php-xdebug.ini
    content: |
      zend_extension=xdebug.so
      xdebug.remote_enable=1
      xdebug.remote_connect_back=0
      xdebug.remote_host=192.168.64.1
      xdebug.remote_port=9000
      xdebug.remote_autostart=1
      xdebug.idekey=PHPSTORM
runcmd:
  - sed -i 's|127.0.0.53|1.1.1.1|g' /etc/resolv.conf
  - add-apt-repository --no-update -y ppa:nginx/stable
  - add-apt-repository --no-update -y ppa:ondrej/php
  - echo "CRAFT_NITRO=1" >> /etc/environment
  - echo "DB_USER=nitro" >> /etc/environment
  - echo "DB_PASSWORD=nitro" >> /etc/environment
  - mkdir -p /home/ubuntu/.composer
  - echo "COMPOSER_HOME=/home/ubuntu/.composer" >> /etc/environment
  - curl -sS https://getcomposer.org/installer -o composer-setup.php
  - php composer-setup.php --install-dir=/usr/local/bin --filename=composer
  - rm composer-setup.php
  - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  - sudo add-apt-repository --no-update -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  - wget -q -O - https://packages.blackfire.io/gpg.key | sudo apt-key add -
  - echo "deb http://packages.blackfire.io/debian any main" | sudo tee /etc/apt/sources.list.d/blackfire.list
  - apt-get update -y
  - apt-get install -y nginx docker-ce docker-ce-cli containerd.io
  - usermod -aG docker ubuntu
  - mkdir -p /home/ubuntu/sites
  - mkdir -p /home/ubuntu/.nitro/databases/imports
  - mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d
  - mkdir -p /home/ubuntu/.nitro/databases/mysql/backups
  - mkdir -p /home/ubuntu/.nitro/databases/postgres/conf.d
  - mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d
  - mkdir -p /home/ubuntu/.nitro/databases/postgres/backups
  - chown -R ubuntu:ubuntu /home/ubuntu/
`
