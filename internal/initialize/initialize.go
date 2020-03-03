package initialize

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

var cloudInit = `#cloud-config
packages:
  - redis
write_files:
  # create the main bootstrap script
  - path: /opt/nitro/bootstrap.sh
    content: |
      #!/usr/bin/env bash
      phpversion="$1"
      database="$2"

      # run the correct script depending on the version of PHP
      if [ "$phpversion" == '7.3' ]; then
        echo "7.3" > /home/ubuntu/.php_version
        sudo bash /opt/nitro/php/php-73.sh
      elif [ "$phpversion" == '7.2' ]; then
        echo "7.2" > /home/ubuntu/.php_version
        sudo bash /opt/nitro/php/php-72.sh
      else
        echo "7.4" > /home/ubuntu/.php_version
        sudo bash /opt/nitro/php/php-74.sh
      fi

      # run the correct script depending on the database
      if [ "$database" == 'postgres' ]; then
        # install postgres
        sudo bash /opt/nitro/postgres/install.sh

        # run the postgres setup
        sudo bash /opt/nitro/postgres/setup.sh
      else
        # install mariadb
        sudo bash /opt/nitro/mariadb/install.sh

        # run the mariadb setup
        sudo bash /opt/nitro/mariadb/setup.sh
      fi

      # install nginx
      sudo bash /opt/nitro/nginx/install.sh

      # remove the default site
      sudo rm /etc/nginx/sites-enabled/default
    permissions: '770'
  - path: /opt/nitro/update.sh
    content: |
      #!/usr/bin/env bash
      sudo apt update && sudo apt upgrade -y
    permissions: '770'
  # create the php install scripts
  - path: /opt/nitro/php/php-74.sh
    content: |
      #!/bin/bash
      apt install -y php7.4 php7.4-mbstring php7.4-cli php7.4-curl php7.4-fpm php7.4-gd php7.4-intl php7.4-json \
      php7.4-mysql php7.4-opcache php7.4-pgsql php7.4-zip php7.4-xml php7.4-xdebug
  - path: /opt/nitro/php/php-73.sh
    content: |
      #!/bin/bash
      apt install -y php7.3 php7.3-mbstring php7.3-cli php7.3-curl php7.3-fpm php7.3-gd php7.3-intl php7.3-json \
      php7.3-mysql php7.3-opcache php7.3-pgsql php7.3-zip php7.3-xml php7.3-xdebug
  - path: /opt/nitro/php/php-72.sh
    content: |
      #!/bin/bash
      apt install -y php7.2 php7.2-mbstring php7.2-cli php7.2-curl php7.2-fpm php7.2-gd php7.2-intl php7.2-json \
      php7.2-mysql php7.2-opcache php7.2-pgsql php7.2-zip php7.2-xml php7.2-xdebug
  # create the composer install script
  - path: /opt/nitro/composer-install.sh
    content: |
      # install composer
      php -r "readfile('http://getcomposer.org/installer');" | sudo php -- --install-dir=/usr/bin/ --filename=composer
    permissions: '770'
  # create mariadb scripts
  - path: /opt/nitro/mariadb/setup.sh
    content: |
      # remove bind from mariadb to allow remote connection
      systemctl stop mariadb
      sed -i 's|bind-address|#bind-address|g' /etc/mysql/mariadb.conf.d/50-server.cnf
      systemctl start mariadb

      # create the database and user
      if [ -f /home/ubuntu/.db_password ]; then
          # read the file contents for the password
          export DB_USER_PASS=$(cat /home/ubuntu/.db_password)
      else
          # create a random password and store it
          export USER_PASS=$(openssl rand -base64 16)
          echo "$USER_PASS" > /home/ubuntu/.db_password
          export DB_USER_PASS=$(cat /home/ubuntu/.db_password)
      fi

      sudo sed -i 's|CHANGEME|'$DB_USER_PASS'|g' /opt/nitro/mariadb/init.sql

      sudo mysql -u root < /opt/nitro/mariadb/init.sql
    permissions: '770'
  - path: /opt/nitro/mariadb/install.sh
    content: |
      #!/bin/bash
      apt install -y mariadb-server
  - path: /opt/nitro/postgres/setup.sh
    content: |
      # allow remote access to postgres
      sed -i 's|#listen_addresses = 'localhost'|listen_addresses = '*'|g' /etc/postgresql/10/main/postgresql.conf
      sed -i 's|127.0.0.1/32|0.0.0.0/0|g' /etc/postgresql/10/main/pg_hba.conf
      sudo service postgresql restart

      # create the database and user
      if [ -f /home/ubuntu/.db_password ]; then
          # read the file contents for the password
          export DB_USER_PASS=$(cat /home/ubuntu/.db_password)
      else
          # create a random password and store it
          export USER_PASS=$(openssl rand -base64 16)
          echo "$USER_PASS" > /home/ubuntu/.db_password
          export DB_USER_PASS=$(cat /home/ubuntu/.db_password)
      fi

      sudo sed -i 's|CHANGEME|'$DB_USER_PASS'|g' /opt/nitro/postgres/init.sql

      sudo -u postgres psql -U ubuntu < /opt/nitro/postgres/init.sql
    permissions: '770'
  - path: /opt/nitro/postgres/install.sh
    content: |
      #!/bin/bash
      apt install -y postgresql postgresql-contrib
  - path: /opt/nitro/postgres/init.sql
    content: |
      CREATE DATABASE craftcms;
      CREATE ROLE craftcms WITH LOGIN PASSWORD 'CHANGEME';
      GRANT SELECT, INSERT, UPDATE, CREATE, DELETE, REFERENCES, CONNECT ON craftcms TO craftcms; 
  # create nginx install scripts
  - path: /opt/nitro/nginx/install.sh
    content: |
      #!/bin/bash
      apt install -y nginx
  - path: /opt/nitro/nginx/template.conf
    content: |
      server {
          listen 80;

          root /var/www/CHANGEPATH/CHANGEPUBLICDIR;

          index index.html index.htm index.php;

          server_name CHANGESERVERNAME;

          location / {
              try_files $uri $uri/ /index.php$is_args$args;
          }

          location ~ \.php$ {
             include snippets/fastcgi-php.conf;
             fastcgi_pass unix:/var/run/php/phpCHANGEPHPVERSION-fpm.sock;
          }
      }
    permissions: '770'
  - path: /opt/nitro/nginx/add-site.sh
    content: |
      #!/usr/bin/env bash
      NEW_HOST_NAME="$1"
      PHP_VERSION="$2"
      PUBLIC_DIR="$3"

      # make the directories
      sudo mkdir /var/www/"$NEW_HOST_NAME"

      # set permissions for www-data to new directory
      sudo chown -R www-data:www-data /var/www/"$NEW_HOST_NAME"

      # copy the nitro nginx template into sites-available/default
      sudo cp /opt/nitro/nginx/template.conf /etc/nginx/sites-available/"$NEW_HOST_NAME"

      # change the variables
      sudo sed -i 's|CHANGEPATH|'$NEW_HOST_NAME'|g' /etc/nginx/sites-available/"$NEW_HOST_NAME"
      sudo sed -i 's|CHANGESERVERNAME|'$NEW_HOST_NAME'|g' /etc/nginx/sites-available/"$NEW_HOST_NAME"
      sudo sed -i 's|CHANGEPUBLICDIR|'$PUBLIC_DIR'|g' /etc/nginx/sites-available/"$NEW_HOST_NAME"
      sudo sed -i 's|CHANGEPHPVERSION|'$PHP_VERSION'|g' /etc/nginx/sites-available/"$NEW_HOST_NAME"

      # enable the nginx site
      sudo ln -s /etc/nginx/sites-available/"$NEW_HOST_NAME" /etc/nginx/sites-enabled/

      # reload nginx
      sudo service nginx reload
    permissions: '0770'
  - path: /opt/nitro/nginx/tail-logs.sh
    content: |
      #!/usr/bin/env bash
      sudo tail -f /var/log/nginx/access.log -f /var/log/nginx/error.log
    permissions: '0770'
runcmd:
  - sudo add-apt-repository -y ppa:nginx/stable
  - sudo add-apt-repository -y ppa:ondrej/php
  - sudo apt update -y
`

// Command it the main entry point when calling the init command to start a new machine
func Command() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize new machine",
		Action: func(c *cli.Context) error {
			return run(c)
		},
		After: func(c *cli.Context) error {
			// if we are bootstrapping, call the command
			if c.Bool("bootstrap") {
				php := c.String("php-version")
				database := c.String("database")
				args := []string{c.App.Name, "--machine", c.String("machine"), "bootstrap", "--php-version", php, "--database", database}
				return c.App.RunContext(c.Context, args)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "bootstrap",
				Usage:       "Bootstrap the machine with defaults",
				Value:       true,
				DefaultText: "true",
			},
			&cli.StringFlag{
				Name:        "php-version",
				Usage:       "Provide version of PHP",
				DefaultText: "7.4",
			},
			&cli.StringFlag{
				Name:        "database",
				Usage:       "Provide version of PHP",
				DefaultText: "mariadb",
			},
			&cli.Int64Flag{
				Name:        "cpus",
				Usage:       "The number of CPUs to assign the machine",
				Required:    false,
				Value:       2,
				DefaultText: "2",
			},
			&cli.StringFlag{
				Name:        "memory",
				Usage:       "The amount of memory to assign the machine",
				Required:    false,
				Value:       "2G",
				DefaultText: "2G",
			},
			&cli.StringFlag{
				Name:        "disk",
				Usage:       "The amount of disk to assign the machine",
				Required:    false,
				Value:       "5G",
				DefaultText: "5G",
			},
		},
	}
}

func run(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	// create the machine
	cmd := exec.Command(
		multipass,
		"launch",
		"--name",
		machine,
		"--cpus",
		strconv.Itoa(c.Int("cpus")),
		"--disk",
		c.String("disk"),
		"--mem",
		c.String("memory"),
		"--cloud-init",
		"-",
	)

	// pass the cloud init as stdin
	cmd.Stdin = strings.NewReader(cloudInit)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
