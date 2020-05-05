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
  - path: /opt/nitro/scripts/site-exists.sh
    content: |
      #!/usr/bin/env bash
      site="$1"
      if test -f /etc/nginx/sites-enabled/"$site"; then
          echo "exists"
      fi
  - path: /opt/nitro/scripts/docker-container-exists.sh
    content: |
      #!/usr/bin/env bash
      NAME="$1"
      if [ -n "$(docker ps -q -f name="$NAME")" ]; then
        echo "exists"
      fi
  - path: /opt/nitro/scripts/docker-exec-import.sh
    content: |
      #!/usr/bin/env bash
      container="$1"
      database="$2"
      filename="$3"
      engine="$4"
      
      if [ "$engine" == "mysql" ]; then
          docker exec -i "$container" mysql -uroot -pnitro -e "CREATE DATABASE IF NOT EXISTS $database;"
          docker exec -i "$container" mysql -uroot -pnitro -e "GRANT ALL ON $database.* TO 'nitro'@'%';"
          docker exec -i "$container" mysql -uroot -pnitro -e "FLUSH PRIVILEGES;"
          cat "$filename" | pv | docker exec -i "$container" mysql -unitro -pnitro "$database" --init-command="SET autocommit=0;"
      else
          docker exec "$container" psql -U nitro -c "CREATE DATABASE $database OWNER nitro;"
          cat "$filename" | pv | docker exec -i "$container" psql -U nitro -d "$database"
      fi
  - path: /opt/nitro/scripts/docker-set-database-user-permissions.sh
    content: |
      #!/usr/bin/env bash
      container="$1"
      engine="$2"

      if [ -z "$container" ]; then
          echo "you must provide a container name"
          exit 1
      fi

      if [ -z "$engine" ]; then
          echo "you must provide a database engine (e.g. mysql or postgres)"
          exit 1
      fi

      if [ "$engine" == "mysql" ]; then
          docker exec "$container" bash -c "while ! mysqladmin ping -h 127.0.0.1 -uroot -pnitro; do echo 'waiting...'; sleep 4; done"
          docker exec "$container" mysql -uroot -pnitro --silent --no-beep -e "GRANT ALL ON *.* TO 'nitro'@'%';"
          docker exec "$container" mysql -uroot -pnitro -e "FLUSH PRIVILEGES;"
          echo "setting root permissions on user nitro"
      else
          docker exec "$container" psql -U postgres -c "ALTER USER nitro WITH SUPERUSER;"
          echo "setting superuser permissions on user nitro"
      fi
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

          location / {
              try_files $uri $uri/ /index.php$is_args$args;
          }

          location ~ \.php$ {
             include snippets/fastcgi-php.conf;
             fastcgi_pass unix:/var/run/php/phpCHANGEPHPVERSION-fpm.sock;
             fastcgi_read_timeout 240;
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
  - curl -sS https://getcomposer.org/installer -o composer-setup.php
  - export COMPOSER_HOME=/home/ubuntu/composer
  - php composer-setup.php --install-dir=/usr/local/bin --filename=composer
  - rm composer-setup.php
  - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  - sudo add-apt-repository --no-update -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  - wget -q -O - https://packages.blackfire.io/gpg.key | sudo apt-key add -
  - echo "deb http://packages.blackfire.io/debian any main" | sudo tee /etc/apt/sources.list.d/blackfire.list
  - apt-get update -y
  - apt install -y nginx docker-ce docker-ce-cli containerd.io
  - usermod -aG docker ubuntu
  - mkdir -p /nitro/sites
  - chown -R ubuntu:ubuntu /nitro/sites
`
