#!/usr/bin/env bash
version="$!"

if [ -z $version ]; then
  echo "you must specify a version of nitro"
  exit 1
fi

if [ $version == "1.0.0-beta.3" ]; then
  echo "running sync script for 1.0.0-beta.3"

  # create setup scripts
  mkdir -p /home/ubuntu/sites
  mkdir -p /home/ubuntu/.nitro/databases/imports
  mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d
  mkdir -p /home/ubuntu/.nitro/databases/mysql/backups
  mkdir -p /home/ubuntu/.nitro/databases/postgres/conf.d
  mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d
  mkdir -p /home/ubuntu/.nitro/databases/postgres/backups
  chown -R ubuntu:ubuntu /home/ubuntu/.nitro
  chown -R ubuntu:ubuntu /home/ubuntu/sites

  # create the files
  cat >"/home/ubuntu/.nitro/databases/mysql/conf.d/mysql.conf" <<-EndOfMessage
[mysqld]
max_allowed_packet=1000M
wait_timeout=3000
EndOfMessage

  cat >"/home/ubuntu/.nitro/databases/mysql/setup.sql" <<-EndOfMessage
CREATE USER IF NOT EXISTS 'nitro'@'localhost' IDENTIFIED BY 'nitro';
CREATE USER IF NOT EXISTS 'nitro'@'%' IDENTIFIED BY 'nitro';
GRANT ALL PRIVILEGES ON *.* TO 'nitro'@'localhost' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'nitro'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
EndOfMessage

  cat >"/home/ubuntu/.nitro/databases/postgres/setup.sql" <<-EndOfMessage
ALTER USER nitro WITH SUPERUSER;
EndOfMessage

  echo "removing old scripts"
  sudo rm -rf /opt/nitro/scripts

  echo "setup script has completed!"
  exit 0
fi
