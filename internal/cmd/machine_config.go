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
write_files:
  - path: /opt/nitro/scripts/site-exists.sh
    content: |
      #!/usr/bin/env bash
      site="$1"
      if test -f /etc/nginx/sites-enabled/"$site"; then
          echo "exists"
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
  - sudo add-apt-repository --no-update -y ppa:nginx/stable
  - sudo add-apt-repository --no-update -y ppa:ondrej/php
  - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  - sudo add-apt-repository --no-update -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  - wget -q -O - https://packages.blackfire.io/gpg.key | sudo apt-key add -
  - echo "deb http://packages.blackfire.io/debian any main" | sudo tee /etc/apt/sources.list.d/blackfire.list
  - sudo apt-get update -y
  - sudo apt install -y nginx docker-ce docker-ce-cli containerd.io
  - sudo usermod -aG docker ubuntu
  - sudo mkdir -p /nitro/sites
  - sudo chown -R ubuntu:ubuntu /nitro/sites
`
