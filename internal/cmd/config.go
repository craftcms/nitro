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
  - unzip
  - mysql-client
  - postgresql-client
write_files:
  - path: /home/ubuntu/.nitro/databases/mysql/conf.d/5/mysql.cnf
    content: |
      [mysqld]
      max_allowed_packet=256M
      wait_timeout=86400
      default-authentication-plugin=mysql_native_password
  - path: /home/ubuntu/.nitro/databases/mysql/conf.d/8/mysql.cnf
    content: |
      [mysqld]
      max_allowed_packet=256M
      wait_timeout=86400
      default-authentication-plugin=mysql_native_password
      [mysqldump]
      column-statistics=0
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
          error_page 404 /index.php$is_args$args;

          # Root directory location handler
          location / {
              try_files $uri/index.html $uri $uri/ /index.php$is_args$args;
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
              add_header Last-Modified $date_gmt;
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
  - sed -i 's|nameserver 127.0.0.53|nameserver 127.0.0.53\nnameserver 1.1.1.1\nnameserver 1.0.0.1\nnameserver 8.8.8.8\nnameserver 8.8.4.4|g' /etc/resolv.conf
  - add-apt-repository --no-update -y ppa:ondrej/php
  - echo "CRAFT_NITRO=1" >> /etc/environment
  - echo "DB_USER=nitro" >> /etc/environment
  - echo "DB_PASSWORD=nitro" >> /etc/environment
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
  - cp /etc/skel/.bashrc /home/ubuntu/.bashrc
  - cp /etc/skel/.profile /home/ubuntu/.profile
  - cp /etc/skel/.bash_logout /home/ubuntu/.bash_logout
  - sed -i 's|#force_color_prompt=yes|force_color_prompt=yes|g' /home/ubuntu/.bashrc
  - chown -R ubuntu:ubuntu /home/ubuntu/
  - wget https://raw.githubusercontent.com/craftcms/nitro/master/nitrod.sh -O /tmp/nitrod.sh
  - bash /tmp/nitrod.sh
`
